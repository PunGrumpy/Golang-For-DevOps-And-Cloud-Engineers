package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/11-http-login-packaged/pkg/api"
)

func main() {
	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	flag.StringVar(&requestURL, "url", "", "URL to access")
	flag.StringVar(&password, "password", "", "Use a password to access our API")

	flag.Parse()

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Invalid URL: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	apiInstace := api.New(api.Options{
		Password: password,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
	})

	res, err := apiInstace.DoGetRequest(parsedURL.String())
	if err != nil {
		if requestErr, ok := err.(api.RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	if res == nil {
		fmt.Printf("No response\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}
