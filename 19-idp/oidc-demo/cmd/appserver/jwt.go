package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/oidc"
	"github.com/golang-jwt/jwt/v4"
)

// gets token from tokenUrl validating token with jwksUrl and returning token & claims
func getTokenFromCode(tokenUrl, jwksUrl, redirectUri, clientID, clientSecret, code string) (*jwt.Token, *jwt.RegisteredClaims, error) {
	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("client_id", clientID)
	values.Add("client_secret", clientSecret)
	values.Add("redirect_uri", redirectUri)
	values.Add("code", code)

	res, err := http.PostForm(tokenUrl, values)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("Error getting token: %s", body)
	}

	var token oidc.Token

	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, nil, fmt.Errorf("Error unmarshalling token: %s", err)
	}

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token.IDToken, claims, func(*jwt.Token) (interface{}, error) {
		return nil, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing token: %s", err)
	}

	return parsedToken, claims, fmt.Errorf("Not implemented")
}
