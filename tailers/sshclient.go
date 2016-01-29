package tailers

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Name   string            `json:"name"`
	Host   string            `json:"host"`
	Port   int               `json:"port"`
	Config *ssh.ClientConfig `json:"-"`
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print("ERROR opening file\n")
		fmt.Println(err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Print("ERROR parsing file\n")
		fmt.Println(err)
		return nil
	}
	return ssh.PublicKeys(key)
}

func (client *SSHClient) RunCommand(cmd *SSHCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSession(); err != nil {
		return err
	}
	defer session.Close()

	if err = client.prepareCommand(session, cmd); err != nil {
		return err
	}

	err = session.Run(cmd.Path)
	return err
}

func (client *SSHClient) RunCommandAndWait(cmd *SSHCommand) ([]byte, error) {
	var (
		session *ssh.Session
		err     error
		baData  []byte
	)

	if session, err = client.newSession(); err != nil {
		return nil, err
	}
	defer session.Close()

	baData, err = session.Output(cmd.Path)
	return baData, err
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

func (client *SSHClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}
