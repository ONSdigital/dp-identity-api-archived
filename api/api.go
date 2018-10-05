//Package provides HTTP Handlers/HandlerFunc's for creating, updating and deleting user identities
package api

import (
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/healthcheck"
	"github.com/gorilla/mux"
)

//New is a constructor function for creating a new instance of the API.
func New(host string, identityService IdentityService, tokenService TokenService, auditor audit.AuditorService) *API {
	return &API{
		Host:            host,
		IdentityService: identityService,
		Tokens:          tokenService,
		auditor:         auditor,
	}
}

//RegisterEndpoints provides a way to register the HandlerFunc's defined in the api package with a mux.Router.
func (api *API) RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/identity", api.CreateIdentityHandler).Methods("POST")
	r.HandleFunc("/identity", api.GetIdentityHandler).Methods("GET")
	r.HandleFunc("/token", api.CreateTokenHandler).Methods("POST")
	r.Path("/healthcheck").HandlerFunc(healthcheck.Do)
}
