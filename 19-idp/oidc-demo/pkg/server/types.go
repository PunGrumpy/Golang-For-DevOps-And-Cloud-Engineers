package server

import (
	"time"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/users"
)

type Config struct {
	Apps      map[string]AppConfig `yaml:"apps"`
	Url       string               `yaml:"url"`
	LoadError error
}
type AppConfig struct {
	ClientID     string   `yaml:"clientID"`
	ClientSecret string   `yaml:"clientSecret"`
	Issuer       string   `yaml:"issuer"`
	RedirectURIs []string `yaml:"redirectURIs"`
}
type LoginRequest struct {
	ClientID     string
	RedirectURI  string
	ResponseType string
	State        string
	Scope        string
	AppConfig    AppConfig
	CodeIssuedAt time.Time
	User         users.User
}
