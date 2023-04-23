package main

import (
	"fmt"
	"io/ioutil"

	"github.com/wardviaene/golang-for-devops-course/ssh-demo"
)

func main() {
	var (
		err error
	)
	serverKeyBytes, err := ioutil.ReadFile("mykey.pem")
	if err != nil {
		fmt.Printf("unable to read server key: %v", err)
	}

	authorizedKeysBytes, err := ioutil.ReadFile("server.pub")
	if err != nil {
		fmt.Printf("unable to read authorized keys: %v", err)
	}

	if err = ssh.StartServer(serverKeyBytes, authorizedKeysBytes); err != nil {
		fmt.Printf("unable to start SSH server: %v", err)
	}
}
