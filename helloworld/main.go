package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./helloworld <arguments>")
		os.Exit(1)
	}
	
	fmt.Printf("hello world\nOS.arg: %v\nArguments: %v\n", args, args[1:])
}