package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/users"
	"github.com/golang-jwt/jwt/v4"
)

func (s *server) userinfo(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		returnError(w, fmt.Errorf("authorization header is missing"))
		return
	}

	authorizationHeader = strings.Replace(authorizationHeader, "Bearer ", "", -1)

	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(authorizationHeader, claims, func(token *jwt.Token) (interface{}, error) {
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(s.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		returnError(w, fmt.Errorf("parse with claims: %w", err))
		return
	}

	found := false
	for _, aud := range claims.Audience {
		if aud == s.Config.Url+"/userinfo" {
			found = true
		}
	}
	if !found {
		returnError(w, fmt.Errorf("audience is invalid: %s", strings.Join(claims.Audience, ", ")))
		return
	}
	if claims.Subject == "" {
		returnError(w, fmt.Errorf("subject is empty"))
		return
	}

	for _, user := range users.GetAllUsers() {
		if user.Sub == claims.Subject {
			out, err := json.Marshal(user)
			if err != nil {
				returnError(w, fmt.Errorf("json marshal: %w", err))
				return
			}
			w.Write(out)
			return
		}
	}

	returnError(w, fmt.Errorf("user not found"))
}
