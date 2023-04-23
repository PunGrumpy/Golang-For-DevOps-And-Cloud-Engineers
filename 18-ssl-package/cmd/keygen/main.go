package main

import (
	"fmt"
	"os"

	// "github.com/wardviaene/golang-for-devops-course/ssh-demo"
	"github.com/PunGrumpy/go-ssh-key/ssh"
)

func main() {
	var (
		privateKey, publicKey []byte
		err                   error
	)
	if privateKey, publicKey, err = ssh.GenerateKey(); err != nil {
		fmt.Printf("unable to generate key pair: %v", err)
	}
	if err = os.WriteFile("mykey.pem", privateKey, 0600); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile("mykey.pub", publicKey, 0644); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
