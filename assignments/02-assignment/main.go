/*
* 	Assignment 2
*	Using what you have learned, run the test-server again. Try to hit http://localhost:8080/ratelimit at a constant rate of 5 requests in a given second.
*	The server will reply what your current rate is, and will prefix "DONE!" if you're hitting 5 requests in a given second.
*	If you go over 5 requests per second, it'll return a 429 HTTP error code for the next 10 seconds.
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type RateLimiter struct {
	mu sync.Mutex
}

func main() {
	arg := os.Args[1]
	if len(arg) < 1 {
		fmt.Println("Please provide a URL")
		os.Exit(1)
	}
	rateLimiterInstance := &RateLimiter{}
	finish := make(chan string)
	go func(rateLimiterInstance *RateLimiter) {
		for i := 0; i < 5; i++ {
			rateLimiterInstance.RateLimit(arg)
		}
		finish <- "DONE!"
	}(rateLimiterInstance)
	fmt.Printf("\nProgram finished: %s\n", <-finish)
}

func (r *RateLimiter) RateLimit(arg string) {
	r.mu.Lock()
	requestToServer(arg)
	r.mu.Unlock()
}

func requestToServer(arg string) {
	if _, err := url.ParseRequestURI(arg); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	response, err := http.Get(arg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Response: %s", body)
}
