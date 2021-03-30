package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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

var clear map[string]func()

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

	hosts := make([]Host, 0)
	for _, h := range sshCfg.Hosts {
		if h.Patterns[0].String() == "*" {
			continue
		}
		host := Host{
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

func connectServer(host Host) error {
	fileData, err := ioutil.ReadFile(host.IdentityFile)
	if err != nil {
		return err
	}
	singer, err := ssh.ParsePrivateKey(fileData)
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.HostName, host.Port), &ssh.ClientConfig{
		User:            host.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(singer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

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

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}
	if err = session.RequestPty(term, h, w, ssh.TerminalModes{
		ssh.ECHO: 1,
	}); err != nil {
		return err
	}

	exist := make(chan struct{})
	defer func() {
		exist <- struct{}{}
	}()
	go func() {
		for {
			select {
			case <-exist:
				break
			default:
				_, _ = session.SendRequest("print", false, nil)
			}
			time.Sleep(time.Second)
		}
	}()

	if err = session.Shell(); err != nil {
		return err
	}
	_ = session.Wait()
	return nil
}
