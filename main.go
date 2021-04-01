package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hnakamur/go-scp"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Host struct {
	Host         string `json:"host"`
	HostName     string `json:"host_name"`
	User         string `json:"user"`
	Port         int    `json:"port"`
	IdentityFile string `json:"identity_file"`
	Comment      string `json:"comment"`
}

type Session struct {
	hostConfig *Host
	session    *ssh.Session
	client     *ssh.Client
	ctx        context.Context
	cancel     context.CancelFunc
	cmd        *cmdEntity

	writeLock   *sync.RWMutex
	stdinPiper  io.WriteCloser
	stdoutPiper io.Reader
}

type cmdEntity struct {
	buf       []byte
	hasTab    bool
	hasUpDown bool
}

const (
	bash        = "-bash: %s: "
	cmdNotFound = "command not found"
	legalWords  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890./_- "
)

var (
	clear          map[string]func()
	allCmd         = []string{"down", "up"}
	allCmdNotFound []string
)

func init() {
	clear = make(map[string]func())
	unixClearFunc := func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
	clear["linux"] = unixClearFunc
	clear["darwin"] = unixClearFunc
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	for _, v := range allCmd {
		allCmdNotFound = append(allCmdNotFound, fmt.Sprintf(bash, v)+cmdNotFound)
	}
}

func main() {
	f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	if err != nil {
		panic(err)
	}
	sshCfg, err := ssh_config.Decode(f)
	if err != nil {
		panic(err)
	}

	hosts := make([]*Host, 0)
	for _, h := range sshCfg.Hosts {
		if h.Patterns[0].String() == "*" {
			continue
		}
		host := &Host{
			Host:         h.Patterns[0].String(),
			User:         os.Getenv("USER"),
			Port:         22,
			IdentityFile: filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"),
			Comment:      h.EOLComment,
		}

		for _, node := range h.Nodes {
			switch v := node.(type) {
			case *ssh_config.KV:
				switch v.Key {
				case "HostName":
					host.HostName = v.Value
				case "Port":
					host.Port, _ = strconv.Atoi(v.Value)
				case "IdentityFile":
					host.IdentityFile = strings.ReplaceAll(v.Value, "~", os.Getenv("HOME"))
				case "User":
					host.User = v.Value
				}
			}
		}
		hosts = append(hosts, host)
	}

	prompt := promptui.Select{
		Size:  20,
		Label: "选择机器",
	}
	var items []string
	for _, host := range hosts {
		items = append(items, host.HostName+" "+host.Comment)
	}
	prompt.Items = items

	for {
		clear[runtime.GOOS]()
		idx, _, err := prompt.Run()
		if err != nil {
			if err.Error() == "^C" {
				return
			}
			panic(err)
		}
		host := hosts[idx]
		if err := connectServer(host); err != nil {
			panic(err)
		}
	}
}

func (h *Host) getClient() (*ssh.Client, error) {
	fileData, err := ioutil.ReadFile(h.IdentityFile)
	if err != nil {
		return nil, err
	}
	singer, err := ssh.ParsePrivateKey(fileData)
	if err != nil {
		return nil, err
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.HostName, h.Port), &ssh.ClientConfig{
		User:            h.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(singer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func connectServer(host *Host) error {
	client, err := host.getClient()
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}
	if err = session.RequestPty(term, termHeight, termWidth, ssh.TerminalModes{
		ssh.ECHO: 1,
	}); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stdinPiper, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer stdinPiper.Close()

	stdoutPiper, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPiper, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, stderrPiper)

	s := &Session{
		hostConfig: host,
		session:    session,
		client:     client,
		ctx:        ctx,
		cmd: &cmdEntity{
			buf: make([]byte, 0, 128),
		},
		writeLock:   &sync.RWMutex{},
		stdinPiper:  stdinPiper,
		stdoutPiper: stdoutPiper,
	}
	go s.watchWinch()
	go s.ping()

	go s.writePiperStdin()
	go s.readPiperStdout()

	if err = session.Shell(); err != nil {
		return err
	}
	_ = session.Wait()
	return nil
}

func (s *Session) watchWinch() error {
	sigwinchCh := make(chan os.Signal, 1)
	signal.Notify(sigwinchCh, syscall.SIGWINCH)

	fd := int(os.Stdin.Fd())
	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case sigwinch := <-sigwinchCh:
			if sigwinch == nil {
				return nil
			}
			currTermWidth, currTermHeight, err := terminal.GetSize(fd)

			if currTermHeight == termHeight && currTermWidth == termWidth {
				continue
			}
			_ = s.session.WindowChange(currTermHeight, currTermWidth)
			if err != nil {
				continue
			}
			termWidth, termHeight = currTermWidth, currTermHeight
		}
	}
}

func (s *Session) ping() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			_, _ = s.session.SendRequest("print", false, nil)
			time.Sleep(time.Second)
		}
	}
}

func (s *Session) writePiperStdin() error {
	fd := int(os.Stdin.Fd())
	buf := make([]byte, 128)
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			n, err := syscall.Read(fd, buf)
			if err != nil {
				return err
			}
			if n > 0 {
				s.writeLock.Lock()
				if _, err := s.stdinPiper.Write(buf[:n]); err != nil {
					s.writeLock.Unlock()
					return err
				}
				s.writeLock.Unlock()

				key := GetKey(buf[:n])
				switch key {
				case Up, Down:
					s.cmd.hasUpDown = true
				case Tab:
					s.cmd.hasTab = true
				case Backspace:
					if len(s.cmd.buf) > 0 {
						s.cmd.buf = s.cmd.buf[:len(s.cmd.buf)-1]
					}
				case Enter, ControlM:
					if len(s.cmd.buf) > 0 {
						for _, cmd := range allCmd {
							if bytes.HasPrefix(s.cmd.buf, []byte(cmd)) {
								if err := s.runCmd(string(s.cmd.buf)); err != nil {
									return err
								}
								break
							}
						}
						s.cmd.buf = make([]byte, 0, 128)
					}
				default:
					for _, b := range buf[:n] {
						if bytes.Contains([]byte(legalWords), []byte{b}) {
							s.cmd.buf = append(s.cmd.buf, b)
						}
					}
				}
			}
		}
	}
}

func (s *Session) readPiperStdout() error {
	buf := make([]byte, 128)

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			n, err := s.stdoutPiper.Read(buf)
			if err != nil {
				return err
			}
			if n > 0 {
				ok, result := isSelfCmd(buf[:n])
				if !ok {
					if _, err := os.Stdout.Write(buf[:n]); err != nil {
						return err
					}
					if s.cmd.hasTab {
						for _, b := range buf[:n] {
							if bytes.Contains([]byte(legalWords), []byte{b}) {
								s.cmd.buf = append(s.cmd.buf, b)
							}
						}
						s.cmd.hasTab = false
					}

					if s.cmd.hasUpDown {
						s.cmd.buf = buf[:n]
						s.cmd.hasUpDown = false
					}
					continue
				}
				if len(result) > 0 {
					if _, err := os.Stdout.Write(result); err != nil {
						return err
					}
				}
			}
		}
	}
}

func (s *Session) runCmd(cmdStr string) error {
	cmdParams := strings.Split(cmdStr, " ")
	if len(cmdParams) == 0 {
		return nil
	}

	cmd := cmdParams[0]
	switch cmd {
	case "down", "up":
		if len(cmdParams) < 2 {
			return nil
		}

		client, err := s.hostConfig.getClient()
		if err != nil {
			return err
		}

		fileName := filepath.Base(strings.TrimSpace(cmdParams[1]))
		localPath := "."
		if len(cmdParams) >= 3 {
			localPath = cmdParams[2]
		}
		localPath = filepath.Join(localPath, fileName)

		go func() {
			defer client.Close()

			if cmd == "down" {
				if err := scp.NewSCP(client).ReceiveFile(cmdParams[1], localPath); err != nil && err != io.EOF {
					s.writeLock.Lock()
					_, _ = s.stdinPiper.Write([]byte{0xd})
					_, _ = os.Stdout.Write([]byte(fmt.Sprintf("\r\rdown %s error: %v   ", fileName, err)))
					_, _ = s.stdinPiper.Write([]byte{0xd})
					s.writeLock.Unlock()
					return
				}
				s.writeLock.Lock()
				_, _ = s.stdinPiper.Write([]byte{0xd})
				_, _ = os.Stdout.Write([]byte(fmt.Sprintf("\r\rdown %s success   ", fileName)))
				_, _ = s.stdinPiper.Write([]byte{0xd})
				s.writeLock.Unlock()
			}

			if cmd == "up" {
				if err := scp.NewSCP(client).SendFile(cmdParams[1], localPath); err != nil && err != io.EOF {
					s.writeLock.Lock()
					_, _ = s.stdinPiper.Write([]byte{0xd})
					_, _ = os.Stdout.Write([]byte(fmt.Sprintf("\r\rup %s error: %v   ", fileName, err)))
					_, _ = s.stdinPiper.Write([]byte{0xd})
					s.writeLock.Unlock()
					return
				}
				s.writeLock.Lock()
				_, _ = s.stdinPiper.Write([]byte{0xd})
				_, _ = os.Stdout.Write([]byte(fmt.Sprintf("\r\rup %s success   ", fileName)))
				_, _ = s.stdinPiper.Write([]byte{0xd})
				s.writeLock.Unlock()
			}
		}()

	default:
		return nil
	}
	return nil
}

func isSelfCmd(cmd []byte) (bool, []byte) {
	cmd = bytes.TrimSpace(cmd)
	for _, v := range allCmdNotFound {
		if bytes.Contains(cmd, []byte(v)) {
			return true, bytes.ReplaceAll(cmd, []byte(v), []byte("\r"))
		}
	}
	return false, nil
}
