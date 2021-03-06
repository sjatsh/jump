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

	"github.com/gogf/gf/util/gconv"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"github.com/sjatsh/go-scp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Host struct {
	hosts                           []string
	Index                           int
	Env                             string
	Host                            string
	AddKeysToAgent                  string
	AddressFamily                   string
	BindAddress                     string
	ChallengeResponseAuthentication string
	Compression                     string
	CompressionLevel                int
	ConnectionAttempts              int
	ConnectTimeout                  int
	ControlMaster                   string
	ControlPath                     string
	ControlPersist                  string
	GatewayPorts                    string
	HostName                        string
	IdentitiesOnly                  string
	IdentityFile                    string
	LocalCommand                    string
	LocalForward                    string
	PasswordAuthentication          string
	PermitLocalCommand              string
	Port                            int
	ProxyCommand                    string
	User                            string
	Comment                         string
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

	idx := 0
	hosts := make([]*Host, 0)
	for _, h := range sshCfg.Hosts {
		if h.Patterns[0].String() == "*" {
			continue
		}
		idx++
		host := &Host{
			Env:          "default",
			Index:        idx,
			Host:         h.Patterns[0].String(),
			User:         os.Getenv("USER"),
			Port:         22,
			IdentityFile: filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"),
			Comment:      h.EOLComment,
		}

		params := make(map[string]string)
		for _, node := range h.Nodes {
			switch v := node.(type) {
			case *ssh_config.KV:
				params[v.Key] = v.Value
			}
		}
		if err := gconv.Struct(params, host); err != nil {
			panic(err)
		}

		hostSlice := strings.Split(host.Host, "_")
		host.hosts = hostSlice
		if len(hostSlice) > 1 {
			host.Env = hostSlice[len(hostSlice)-1]
		}
		host.IdentityFile = strings.ReplaceAll(host.IdentityFile, "~", os.Getenv("HOME"))
		hosts = append(hosts, host)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   "\U0001F449 {{ .Index | cyan }}: {{ .Env | cyan }} {{ .User | green }} {{ .HostName | yellow }} {{ .Comment | white }}",
		Inactive: "  {{ .Index | cyan }}: {{ .Env | cyan }} {{ .User | green }} {{ .HostName | yellow }} {{ .Comment | white }}",
		Selected: "\U0001F449 {{ .Index | cyan }}: {{ .Env | cyan }} {{ .User | green }} {{ .HostName | yellow }} {{ .Comment | white }}",
	}

	searcher := func(input string, index int) bool {
		host := hosts[index]
		number, errNumber := strconv.Atoi(input)
		if errNumber == nil && number == host.Index {
			return true
		}
		if strings.Contains(host.Env, input) {
			return true
		}
		if strings.Contains(host.User, input) {
			return true
		}
		if strings.Contains(host.HostName, input) {
			return true
		}
		if strings.Contains(host.Comment, input) {
			return true
		}
		for _, h := range host.hosts {
			if strings.Contains(h, input) {
				return true
			}
		}
		return false
	}

	prompt := promptui.Select{
		Size:              20,
		Label:             "????????????",
		Items:             hosts,
		Templates:         templates,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	for {
		clear[runtime.GOOS]()
		idx, _, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
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
				case KeyUp, KeyDown:
					s.cmd.hasUpDown = true
				case KeyTab:
					s.cmd.hasTab = true
				case KeyBackspace:
					if len(s.cmd.buf) > 0 {
						s.cmd.buf = s.cmd.buf[:len(s.cmd.buf)-1]
					}
				case KeyControlC:
					s.cmd.buf = make([]byte, 0, 128)
				case KeyEnter, KeyControlM:
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
					s.cmd.buf = append(s.cmd.buf, buf[:n]...)
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
						s.cmd.buf = append(s.cmd.buf, buf[:n]...)
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
		scpClient := scp.NewSCP(client)

		fileName := filepath.Base(strings.TrimSpace(cmdParams[1]))
		localPath := "."
		if len(cmdParams) >= 3 {
			localPath = cmdParams[2]
		}
		localPath = filepath.Join(localPath, fileName)

		go func() {
			defer client.Close()

			switch cmd {
			case "down":
				if err := scpClient.ReceiveFile(cmdParams[1], localPath); err != nil && err != io.EOF {
					_ = s.sendMsg(fmt.Sprintf("\r\rdown %s error: %v   ", fileName, err))
					return
				}
				_ = s.sendMsg(fmt.Sprintf("\r\rdown %s success   ", fileName))
			case "up":
				if err := scpClient.SendFile(cmdParams[1], localPath); err != nil && err != io.EOF {
					_ = s.sendMsg(fmt.Sprintf("\r\rup %s error: %v   ", fileName, err))
					return
				}
				_ = s.sendMsg(fmt.Sprintf("\r\rup %s success   ", fileName))
			}
		}()

	default:
		return nil
	}
	return nil
}

func (s *Session) sendMsg(msg string) error {
	s.writeLock.Lock()
	defer s.writeLock.Unlock()
	if _, err := s.stdinPiper.Write(GetCode(KeyEnter)); err != nil {
		return err
	}
	if _, err := os.Stdout.Write([]byte(msg)); err != nil {
		return err
	}
	if _, err := s.stdinPiper.Write(GetCode(KeyEnter)); err != nil {
		return err
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
