package server

import (
	"fmt"
	"net/http"
	"strings"
)

type server struct {
	PrivateKey   []byte
	Config       Config
	LoginRequest map[string]LoginRequest
}

func newServer(privateKey []byte, config Config) *server {
	return &server{
		PrivateKey:   privateKey,
		Config:       config,
		LoginRequest: make(map[string]LoginRequest),
	}
}

func Start(httpServer *http.Server, privateKey []byte, config Config) error {
	s := newServer(privateKey, config)

	fmt.Printf("Starting server on port %s\n", strings.Split(httpServer.Addr, ":")[1])

	http.HandleFunc("/authorization", s.authorization)
	http.HandleFunc("/token", s.token)
	http.HandleFunc("/login", s.login)
	http.HandleFunc("/jwks.json", s.jwks)
	http.HandleFunc("/.well-known/openid-configuration", s.discovery)
	http.HandleFunc("/userinfo", s.userinfo)

	return httpServer.ListenAndServe()
}

func returnError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))
	fmt.Printf("Error: %s\n", err.Error())
}
