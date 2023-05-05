package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/oidc"
	"github.com/golang-jwt/jwt/v4"
)

func (s *server) token(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		returnError(w, fmt.Errorf("method not allowed (it should be POST)"))
		return
	}
	if err := r.ParseForm(); err != nil {
		returnError(w, fmt.Errorf("failed to parse form: %w", err))
		return
	}
	if r.PostForm.Get("grant_type") != "authorization_code" {
		returnError(w, fmt.Errorf("grant_type should be authorization_code: %s", r.PostForm.Get("grant_type")))
		return
	}
	loginRequest, ok := s.Codes[r.PostForm.Get("code")]
	if !ok {
		returnError(w, fmt.Errorf("invalid code"))
		return
	}
	if time.Now().After(loginRequest.CodeIssuedAt.Add(time.Minute * 10)) {
		returnError(w, fmt.Errorf("code expired"))
		return
	}
	if loginRequest.ClientID != r.PostForm.Get("client_id") {
		returnError(w, fmt.Errorf("invalid client_id"))
		return
	}
	if loginRequest.AppConfig.ClientSecret != r.PostForm.Get("client_secret") {
		returnError(w, fmt.Errorf("invalid client_secret"))
		return
	}
	if loginRequest.RedirectURI != r.PostForm.Get("redirect_uri") {
		returnError(w, fmt.Errorf("invalid redirect_uri"))
		return
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(s.PrivateKey)
	if err != nil {
		returnError(w, fmt.Errorf("failed to parse private key: %w", err))
		return
	}
	claims := jwt.MapClaims{
		"iss": s.Config.Url,
		"sub": loginRequest.User.Sub,
		"aud": loginRequest.ClientID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "0-0-0-1" // token.Header["kid"] is "key id"

	signedIDToken, err := token.SignedString(privateKey)
	if err != nil {
		returnError(w, fmt.Errorf("failed to sign token: %w", err))
		return
	}

	// access token
	claims = jwt.MapClaims{
		"iss": s.Config.Url,
		"sub": loginRequest.User.Sub,
		"aud": []string{s.Config.Url + "/userinfo"},
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "0-0-0-1" // token.Header["kid"] is "key id"

	signedAccessToken, err := token.SignedString(privateKey)
	if err != nil {
		returnError(w, fmt.Errorf("failed to sign token: %w", err))
		return
	}

	tokenOutput := oidc.Token{
		AccessToken: signedAccessToken,
		IDToken:     signedIDToken,
		TokenType:   "Bearer",
		ExpiresIn:   60,
	}

	delete(s.LoginRequest, r.PostForm.Get("code"))

	out, err := json.Marshal(tokenOutput)
	if err != nil {
		returnError(w, fmt.Errorf("failed to marshal token: %w", err))
		return
	}
	w.Write(out)
}
