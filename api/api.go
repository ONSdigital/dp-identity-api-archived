package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

//New is a constructor function for creating a new instance of the API.
func New(identityService IdentityService, auditor audit.AuditorService) *API {
	return &API{
		IdentityService: identityService,
		auditor:         auditor,
	}
}

//RegisterEndpoints provides a way to register the HandlerFunc's defined in the api package with a mux.Router.
func (api *API) RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/identity", api.CreateIdentityHandler).Methods("POST")
	r.Path("/healthcheck").HandlerFunc(healthcheck.Do)
}

//CreateIdentityHandler is a POST HTTP handler for creating a new Identity. A request to this endpoint will create an
// audit event showing an attempt to create a new identity was made followed by another event - successful or unsuccessful
// depending on outcome of processing the request.
func (api *API) CreateIdentityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if auditErr := api.auditor.Record(ctx, createIdentityAction, audit.Attempted, nil); auditErr != nil {
		http.Error(w, auditErr.Error(), http.StatusInternalServerError)
		return
	}

	err := api.createIdentity(ctx, r)

	if err != nil {
		api.auditor.Record(ctx, createIdentityAction, audit.Unsuccessful, nil)
		writeErrorResponse(err, w)
		return
	}

	err = api.auditor.Record(ctx, createIdentityAction, audit.Successful, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.InfoCtx(ctx, "createIdentity: identity created successfully", nil)
}

func (api *API) createIdentity(ctx context.Context, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrFailedToReadRequestBody
	}
	defer r.Body.Close()

	var i identity.Model
	if err := json.Unmarshal(body, &i); err != nil {
		return ErrFailedToUnmarshalRequestBody
	}

	return api.IdentityService.Create(ctx, nil)
}

//writeErrorResponse writes a HTTP error back to the response writer. If the err can be cast to apiError then the values
// of err.GetMessage() and err.GetStatus() will be used to set the response body and status code respectively otherwise
// a default 500 status is used with err.Error() for the response body.
func writeErrorResponse(err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError

	if val, ok := errorStatusMapping[err]; ok {
		status = val
	}

	http.Error(w, err.Error(), status)
}
