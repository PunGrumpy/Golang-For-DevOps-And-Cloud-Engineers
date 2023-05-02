package server

import (
	"fmt"
	"net/http"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/oidc"
)

func (s *server) authorization(w http.ResponseWriter, r *http.Request) {
	var (
		clientID     string
		redirectURI  string
		scope        string
		responseType string
		state        string
	)

	if clientID = r.URL.Query().Get("client_id"); clientID == "" {
		returnError(w, fmt.Errorf("client_id is required"))
		return
	}

	if redirectURI = r.URL.Query().Get("redirect_uri"); redirectURI == "" {
		returnError(w, fmt.Errorf("redirect_uri is required"))
		return
	}

	if scope = r.URL.Query().Get("scope"); scope == "" {
		returnError(w, fmt.Errorf("scope is required"))
		return
	}

	if responseType = r.URL.Query().Get("response_type"); responseType == "" {
		returnError(w, fmt.Errorf("response_type is required"))
		return
	}

	if state = r.URL.Query().Get("state"); state == "" {
		returnError(w, fmt.Errorf("state is required"))
		return
	}

	appConfig := AppConfig{}
	for _, app := range s.Config.Apps {
		if app.ClientID == clientID {
			appConfig = app
		}
	}
	if appConfig.ClientID == "" {
		returnError(w, fmt.Errorf("client_id is not allowed/found"))
		return
	}

	found := false
	for _, redirectURIConfig := range appConfig.RedirectURIs {
		if redirectURIConfig == redirectURI {
			found = true
		}
	}
	if !found {
		returnError(w, fmt.Errorf("redirect_uri is not allowed"))
		return
	}

	sessionID, err := oidc.GetRandomString(128)
	if err != nil {
		returnError(w, err)
		return
	}

	s.LoginRequest[sessionID] = LoginRequest{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		Scope:        scope,
		ResponseType: responseType,
		State:        state,
		AppConfig:    appConfig,
	}

	w.Header().Add("location", fmt.Sprint("/login?sessionID="+sessionID))
	w.WriteHeader(http.StatusFound)
}
