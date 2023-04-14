package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("one\n")
	c := make(chan bool) // chan is a keyword that creates a channel
	go testFunction(c)   // go is a keyword that starts a new goroutine
	// goroutines: cpu threads that are managed by the go runtime scheduler
	fmt.Printf("two\n")
	areWeDone := <-c
	fmt.Printf("areWeDone: %v", areWeDone)
}

func testFunction(c chan bool) {
	for i := 0; i < 5; i++ {
		fmt.Printf("checking...\n")
		time.Sleep(1 * time.Second)
	}
	c <- true
}
