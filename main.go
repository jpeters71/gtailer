package main

import (
    "bufio"
	"fmt"
	"log"
	"net/http"
    
    "bozosonparade/gtailer/ws"
    "bozosonparade/gtailer/tailers"
    "github.com/howeyc/gopass"
    "golang.org/x/crypto/ssh"    
	"io"
    "os"

)

func main() {
	log.SetFlags(log.Lshortfile)
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Printf("Please enter the ssh user to use:");
    sshUser, _ := reader.ReadString('\n')
    fmt.Printf("Password:");
    sshPwd := string(gopass.GetPasswdMasked());

	// websocket server
	server := ws.NewServer("/entry")
	go server.Listen()

    go startSshClient(server, sshUser, sshPwd, "phxedupub11.qa", "phxedupub11.qa.aptimus.net")
    go startSshClient(server, sshUser, sshPwd, "phxedupub12.qa", "phxedupub12.qa.aptimus.net")

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))
    
	//log.Fatal(http.ListenAndServe(":8080", nil))
    log.Fatal(http.ListenAndServeTLS(":7443", "cert.pem", "key.pem", nil))
}

func startSshClient(server *ws.Server, sshUser string, sshPwd string, hostName string, hostAddr string) {
    sshConfig := &ssh.ClientConfig {
        User: sshUser,
        Auth: [] ssh.AuthMethod {
            //tailers.PublicKeyFile("c:/jetbrains/gohome/bozo2-pair.pem"),
            ssh.Password(sshPwd),
            
        },
    }
    client := &tailers.SSHClient {
        Name: hostName,
        Host: hostAddr,
        Port: 22,
        Config: sshConfig,
    }

    //var buff bytes.Buffer;
    r,w := io.Pipe()
    
    go func() {
        reader := bufio.NewReader(r)
        for {
            baLine, _, _ := reader.ReadLine();
            strLine := string(baLine[:])
            if len(strLine) > 0 {
                log.Printf("BOZO: %s\n", strLine)
                msg := ws.Message {client.Host, client.Name, strLine}
                server.SendAll(&msg)
            } 
        }    
    }()
    
    
	cmd := &tailers.SSHCommand{
		//Path:   "ls -lat /home",
		//Path:   "tail -f /var/log/boot.log /var/log/auth.log /var/log/kern.log /home/ubuntu/bozo.tst",
        Path:   "tail -f /cust/appserver/logs/*.out /cust/appserver/logs/*.log /cust/aem/crx-quickstart/logs/*.log",
		Env:    []string{},
		Stdin:  os.Stdin,
		Stdout: w,
		Stderr: os.Stderr,
	}    
    
    log.Printf("Running command: %s\n", cmd.Path)
	if err := client.RunCommand(cmd); err != nil {
		log.Fatal("command run error: %s\n", err)
		os.Exit(1)
	}
    log.Printf("Done: \n")
    
}