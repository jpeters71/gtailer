package main

import (
    "bufio"
	"log"
	"net/http"
    
    "bozosonparade/gtailer/ws"
    "bozosonparade/gtailer/tailers"
    "golang.org/x/crypto/ssh"    
	"io"
    "os"

)

func main() {
	log.SetFlags(log.Lshortfile)

	// websocket server
	server := ws.NewServer("/entry")
	go server.Listen()

    go startSshClient(server)

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startSshClient(server *ws.Server) {
    sshConfig := &ssh.ClientConfig {
        User: "ubuntu",
        Auth: [] ssh.AuthMethod {
            tailers.PublicKeyFile("F:/tools/pems/bozo-pair.pem"),
        },
    }
    client := &tailers.SSHClient {
        Name: "bozo-test-aws",
        Host: "ec2-52-35-140-237.us-west-2.compute.amazonaws.com",
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
		Path:   "tail -f /var/log/boot.log /var/log/auth.log /var/log/kern.log /home/ubuntu/bozo.tst",
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