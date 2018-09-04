package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/healthcheck"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

//New is a constructor function for creating a new instance of the API.
func New(host string, identityService IdentityService, auditor audit.AuditorService) *API {
	return &API{
		Host:            host,
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
// depending on outcome of processing the request.If a  request is successful then a URL to the new identity will be
// returned as a HTTP location header in the response
func (api *API) CreateIdentityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if auditErr := api.auditor.Record(ctx, createIdentityAction, audit.Attempted, nil); auditErr != nil {
		http.Error(w, auditErr.Error(), http.StatusInternalServerError)
		return
	}

	response, err := api.createIdentity(ctx, r)

	if err != nil {
		api.auditor.Record(ctx, createIdentityAction, audit.Unsuccessful, nil)
		writeErrorResponse(ctx, err, w)
		return
	}

	err = api.auditor.Record(ctx, createIdentityAction, audit.Successful, common.Params{"id": response.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	api.writeJSONBody(ctx, w, response, http.StatusCreated)
	log.InfoCtx(ctx, "createIdentity: identity created successfully", log.Data{"id": response.ID})
}

func (api *API) createIdentity(ctx context.Context, r *http.Request) (*IdentityCreated, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrFailedToReadRequestBody
	}
	defer r.Body.Close()

	var i identity.Model
	if err := json.Unmarshal(body, &i); err != nil {
		return nil, ErrFailedToUnmarshalRequestBody
	}

	id, err := api.IdentityService.Create(ctx, &i)
	if err != nil {
		return nil, err
	}

	return &IdentityCreated{
		URI: fmt.Sprintf(identityURIFormat, api.Host, id),
		ID:  id,
	}, nil
}

func (api *API) writeJSONBody(ctx context.Context, w http.ResponseWriter, i interface{}, status int) {
	b, err := json.Marshal(i)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "failed to marshal object to JSON"), log.Data{"object": i})
		writeErrorResponse(ctx, err, w)
		return
	}
	w.Header().Set(headerContentType, mimeTypeJSON)
	w.WriteHeader(status)
	w.Write(b)
}

//writeErrorResponse writes a HTTP error back to the response writer. If the err can be cast to apiError then the values
// of err.GetMessage() and err.GetStatus() will be used to set the response body and status code respectively otherwise
// a default 500 status is used with err.Error() for the response body.
func writeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError

	if val, ok := errorStatusMapping[err]; ok {
		status = val
	}

	log.ErrorCtx(ctx, errors.Wrap(err, "writing error response"), log.Data{"status": status})
	http.Error(w, err.Error(), status)
}
