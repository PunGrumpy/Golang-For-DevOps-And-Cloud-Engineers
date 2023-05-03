package server

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/oidc"
	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/users"
)

//go:embed templates/*
var templateFs embed.FS

func (s *server) login(w http.ResponseWriter, r *http.Request) {

	// POST method
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			returnError(w, fmt.Errorf("parse form error: %w", err))
		}
		sessionID := r.PostForm.Get("sessionID")
		loginRequest, ok := s.LoginRequest[sessionID]
		if !ok {
			returnError(w, fmt.Errorf("no sessionID provided"))
			return
		}
		auth, user, err := users.Auth(r.PostForm.Get("login"), r.PostForm.Get("password"), "")
		if err != nil {
			returnError(w, fmt.Errorf("authenication error: %w", err))
		}

		if auth {
			code, err := oidc.GetRandomString(64)
			if err != nil {
				returnError(w, fmt.Errorf("random string error: %w", err))
			}
			loginRequest.CodeIssuedAt = time.Now()
			loginRequest.User = user
			s.Codes[code] = loginRequest

			delete(s.LoginRequest, sessionID)

			w.Header().Add("location", fmt.Sprintf("%s?code=%s&state=%s", loginRequest.RedirectURI, code, loginRequest.State))
			w.WriteHeader(r.Response.StatusCode)
		} else {
			w.Write([]byte("login failed"))
			w.WriteHeader(r.Response.StatusCode)
		}
		return
	}

	// GET method
	var (
		sessionID string
	)
	if sessionID = r.URL.Query().Get("sessionID"); sessionID == "" {
		returnError(w, fmt.Errorf("no sessionID provided"))
		return
	}
	// to access the login template:
	templateFile, err := templateFs.Open("templates/login.html")
	if err != nil {
		returnError(w, fmt.Errorf("error opening template file: %w", err))
		return
	}
	templateFileBytes, err := io.ReadAll(templateFile)
	if err != nil {
		returnError(w, fmt.Errorf("error reading template file: %w", err))
	}

	templateFileStr := strings.Replace(string(templateFileBytes), "$SESSIONID", sessionID, -1)
	w.Write([]byte(templateFileStr))
}
