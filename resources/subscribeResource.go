package resources

import (
	"bozosonparade/gsh"
	"bozosonparade/gtailer/tailers"
	"bozosonparade/gtailer/ws"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/ssh"
)

var SshUser string
var SshPwd string

var hostOper = regexp.MustCompile(`/subscribe/([^/]+)/([^/]+)`)

var server *ws.Server

// SubscribeResourceHandler handles requests to the /subscribe/ path
func SubscribeResourceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PATH: %s\n", r.URL.Path)
	aaMatches := hostOper.FindAllStringSubmatch(r.URL.Path, -1)

	if aaMatches == nil || len(aaMatches) == 0 || len(aaMatches[0]) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid path specified"))
		return
	}
	strHost := aaMatches[0][1]
	strOperation := aaMatches[0][2]
	fmt.Printf("host: %s, oper: %s\n", strHost, strOperation)

	host := gsh.CurrentConfig.GetHost(strHost)
	if host == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid host"))
		return
	}
	operation := gsh.CurrentConfig.GetOperation(strOperation)
	if operation == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid operation"))
		return
	}

	// Make sure the host supports this operation
	if !host.SupportsOp(strOperation) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid operation"))
		return
	}
	strHostAddr := host.Host
	if !strings.HasSuffix(strHostAddr, gsh.CurrentConfig.DefaultSuffix) {
		strHostAddr += "." + gsh.CurrentConfig.DefaultSuffix
	}
	if operation.IsStreaming {
		// websocket server
		if server == nil {
			server = ws.NewServer("/entry")
			go server.Listen()
		}
		go startSSHClient(server, SshUser, SshPwd, host.Name, strHostAddr, operation)
	} else {
		runCommand(w, host.Name, strHostAddr, operation)
	}
}

func runCommand(w http.ResponseWriter, strHostName string, strHostAddr string, op *gsh.Operation) {
	sshConfig := &ssh.ClientConfig{
		User: SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(SshPwd),
		},
	}
	client := &tailers.SSHClient{
		Name:   strHostName,
		Host:   strHostAddr,
		Port:   22,
		Config: sshConfig,
	}

	buf := new(bytes.Buffer)
	cmd := &tailers.SSHCommand{
		Path:   op.ShellCmd,
		Env:    []string{},
		Stdin:  os.Stdin,
		Stdout: buf,
		Stderr: os.Stdout,
	}

	log.Printf("Running command: %s\n", cmd.Path)

	var aBuf []byte
	var err error
	if aBuf, err = client.RunCommandAndWait(cmd); err != nil {
		log.Fatalf("command run error: %s\n", err)
		os.Exit(1)
	}
	//log.Printf("\nBOASO\n%s\n\n", string(aBuf))
	w.Write(aBuf)
}

func startSSHClient(server *ws.Server, sshUser string, sshPwd string, hostName string, hostAddr string, op *gsh.Operation) {
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPwd),
		},
	}
	client := &tailers.SSHClient{
		Name:   hostName,
		Host:   hostAddr,
		Port:   22,
		Config: sshConfig,
	}

	//var buff bytes.Buffer;
	r, w := io.Pipe()

	go func() {
		reader := bufio.NewReader(r)
		for {
			baLine, _, _ := reader.ReadLine()
			strLine := string(baLine[:])
			if len(strLine) > 0 {
				log.Printf("RECVD: %s\n", strLine)
				msg := ws.Message{
					Host: client.Host,
					Name: client.Name,
					Text: strLine}
				server.SendAll(&msg)
			}
		}
	}()

	cmd := &tailers.SSHCommand{
		Path:   op.ShellCmd,
		Env:    []string{},
		Stdin:  os.Stdin,
		Stdout: w,
		Stderr: os.Stderr,
	}

	log.Printf("Running command: %s\n", cmd.Path)
	if err := client.RunCommand(cmd); err != nil {
		log.Fatalf("command run error: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Done: \n")

}
