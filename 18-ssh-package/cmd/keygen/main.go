package main

import (
	"fmt"
	"os"

	ssh "github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/18-ssh-package"
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
