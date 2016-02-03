package main

// This app is a server for client requests for ssh type stuff.

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/howeyc/gopass"

	"bozosonparade/gsh"
	"bozosonparade/gtailer/resources"
	"os"
)

// main is the entry point for this app.
func main() {
	log.SetFlags(log.Lshortfile)
	// Command line flags
	configPtr := flag.String("config", "", "Config name to use.")
	userPtr := flag.String("user", "", "Username to use for sshing.")
	aConfigs := gsh.LoadConfigs()
	strConfig := ""
	flag.Parse()

	// Check command line options
	if len(*configPtr) > 0 {
		strConfig = *configPtr
		fmt.Printf("Using config: %s\n", strConfig)
	} else {
		strConfs := ""
		for _, conf := range aConfigs {
			if len(strConfs) > 0 {
				strConfs += ", "
			}
			strConfs += conf.Name
		}
		strConfig = readLine(fmt.Sprintf("Select configuration to load (%s): ", strConfs))
	}
	for _, conf := range aConfigs {
		if strings.EqualFold(conf.Name, strConfig) {
			gsh.CurrentConfig = &conf
			break
		}
	}
	if gsh.CurrentConfig == nil {
		log.Fatalf("Unable to load config %s", strConfig)
	}
	if len(*userPtr) > 0 {
		resources.SshUser = *userPtr
		fmt.Printf("Using user: %s\n", resources.SshUser)
	} else {
		resources.SshUser = readLine("Please enter the ssh user to use: ")
	}
	fmt.Printf("Password:")
	resources.SshPwd = string(gopass.GetPasswdMasked())

	// other service endpoints
	http.HandleFunc("/hosts", resources.HostsResourceHandler)
	http.HandleFunc("/operations", resources.OperationsResourceHandler)
	http.HandleFunc("/subscribe/", resources.SubscribeResourceHandler)
	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	//log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("gSSH Utilities started and serving data.")
	log.Fatal(http.ListenAndServeTLS(":7443", "cert.pem", "key.pem", nil))
}

func readLine(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)
	text, _ := reader.ReadString('\n')
	// Clear out white space
	text = strings.TrimSpace(text)

	return text
}
