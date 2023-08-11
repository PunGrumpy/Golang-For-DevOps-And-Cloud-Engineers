package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PunGrumpy/Golang-For-DevOps-And-Cloud-Engineers/19-idp/oidc-demo/pkg/oidc"
)

func (s *server) discovery(w http.ResponseWriter, r *http.Request) {
	discovery := oidc.Discovery{
		Issuer:                            s.Config.Url,
		AuthorizationEndpoint:             s.Config.Url + "/authorization",
		TokenEndpoint:                     s.Config.Url + "/token",
		UserinfoEndpoint:                  s.Config.Url + "/userinfo",
		JwksURI:                           s.Config.Url + "/jwks.json",
		ScopesSupported:                   []string{"openid"},
		ResponseTypesSupported:            []string{"code"},
		TokenEndpointAuthMethodsSupported: []string{"none"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		ClaimsSupported:                   []string{"iss", "sub", "aud", "exp", "nbf", "iat"},
		SubjectTypesSupported:             []string{"public"},
	}
	out, err := json.Marshal(discovery)
	if err != nil {
		returnError(w, fmt.Errorf("error marshalling discovery: %s", err.Error()))
		return
	}
	w.Write(out)
}
