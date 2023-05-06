package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/oidc"
	"github.com/joho/godotenv"
)

const redirectUri = "http://localhost:8081/callback"

type app struct {
	states map[string]bool
}

func main() {

	a := app{
		states: make(map[string]bool),
	}

	http.HandleFunc("/", a.index)
	http.HandleFunc("/callback", a.callback)

	fmt.Printf("Server started on port 8081\n")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("ListenAndServe error: %s\n", err)
	}
}

func (a *app) index(w http.ResponseWriter, r *http.Request) {
	loadEnv()
	oidcEndpoint := os.Getenv("OIDC_ENDPOINT")
	if oidcEndpoint == "" {
		returnError(w, fmt.Errorf("OIDC_ENDPOINT is required"))
		return
	}

	discovery, err := oidc.ParseDiscovery(oidcEndpoint + "/.well-known/openid-configuration")
	if err != nil {
		returnError(w, fmt.Errorf("Error parsing discovery: %s", err))
		return
	}

	state, err := oidc.GetRandomString(64)
	if err != nil {
		returnError(w, fmt.Errorf("Error generating state: %s", err))
		return
	}

	a.states[state] = true

	authorizationURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=openid&response_type=code&state=%s", discovery.AuthorizationEndpoint, os.Getenv("CLIENT_ID"), redirectUri, state)
	w.Write([]byte(`
	<!DOCTYPE html>
	<html>
		<head>
			<title>OIDC Demo</title>
			<style>
				body {
				font-family: Arial, sans-serif;
				background-image: url('https://pbs.twimg.com/media/Fsbdo5AWAAEZ-Zq?format=jpg&name=4096x4096');
				background-size: cover;
				background-position: center;
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
				padding: 0;
				}

				.main-form {
					width: auto;
					max-width: 480px;
					margin: 0 auto;
					padding: 10rem;
					background: rgba(0, 0, 0, 0.5);
					backdrop-filter: blur(10px);
					border-radius: 0.5rem;
					box-shadow: 0px 30px 60px rgba(0, 0, 0, 0.1),
					  0px 30px 60px rgba(0, 0, 0, 0.5);
				}

				a {
					margin: 0 auto;
					background-color: #007bff;
					color: #ffffff;
					text-decoration: none;
					padding: 15px 30px;
					border-radius: 0.5rem;
					font-size: 18px;
					transition: cubic-bezier(0.175, 0.885, 0.32, 1.275) 0.5s;
				}

				a:hover {
				background-color: #0056b3;
				}
			</style>
		</head>
		<body>
			<div class="main-form">
				<a href="` + authorizationURL + `">Login with OIDC</a>
			</div>
		</body>
	</html>
	`))
}

func (a *app) callback(w http.ResponseWriter, r *http.Request) {
	loadEnv()
	oidcEndpoint := os.Getenv("OIDC_ENDPOINT")
	if oidcEndpoint == "" {
		returnError(w, fmt.Errorf("OIDC_ENDPOINT is required"))
		return
	}

	discovery, err := oidc.ParseDiscovery(oidcEndpoint + "/.well-known/openid-configuration")
	if err != nil {
		returnError(w, fmt.Errorf("Error parsing discovery: %s", err))
		return
	}

	if _, ok := a.states[r.URL.Query().Get("state")]; !ok {
		returnError(w, fmt.Errorf("Invalid state"))
		return
	}

	_, claims, err := getTokenFromCode(discovery.TokenEndpoint, discovery.JwksURI, redirectUri, os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), r.URL.Query().Get("code"))
	if err != nil {
		returnError(w, fmt.Errorf("Error getting token from code: %s", err))
		return
	}

	w.Write([]byte(fmt.Sprint("Token: ", claims.Subject)))
}

func returnError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
	fmt.Printf("Error: %s\n", err)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %s\n", err)
	}
}
