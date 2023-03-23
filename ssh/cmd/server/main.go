package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jbrcoleman/golang-devops/ssh"
)

func main() {

	var (
		err error
	)
	serverKeyBytes, err := ioutil.ReadFile("server.pem")
	if err != nil {
		log.Fatalf("Failed to load server keys, err:%v", err)
	}

	authorizedKeysBytes, err := ioutil.ReadFile("mykey.pub")
	if err != nil {
		log.Fatalf("Failed to load authorized keys, err:%v", err)
	}

	if err = ssh.StartServer(serverKeyBytes, authorizedKeysBytes); err != nil {
		fmt.Printf("Error: %s \n")
		os.Exit(1)
	}
}
