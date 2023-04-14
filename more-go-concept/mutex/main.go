package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type myType struct {
	conuter int
	mu      sync.Mutex
}

func (m *myType) IncreaseCounter() {
	m.mu.Lock()
	m.conuter++
	m.mu.Unlock()
}

func main() {
	myTypeInstance := myType{}
	finish := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(myTypeInstance *myType) {
			myTypeInstance.mu.Lock()
			fmt.Printf("Input: %d\n", myTypeInstance.conuter)
			myTypeInstance.conuter++
			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			if myTypeInstance.conuter == 5 {
				fmt.Printf("Found counter == 5\n")
			}
			fmt.Printf("Output: %d\n", myTypeInstance.conuter)
			finish <- true
			myTypeInstance.mu.Unlock()
		}(&myTypeInstance)
	}
	for i := 0; i < 10; i++ {
		<-finish
	}
	fmt.Printf("Counter: %d\n", myTypeInstance.conuter)
}
