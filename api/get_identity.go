package api

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"net/http"
)

// GetIdentityHandler is a GET HTTP handler for retrieving an Identity using a token provided in the request header.
// A request to this endpoint will create audit event showing an attempt to get an identity was made followed by another
// event - successful or unsuccessful depending on outcome of processing the request. If a request is successful the
// retrieved identity will be returned in the response.
func (api *API) GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if auditErr := api.auditor.Record(ctx, getIdentityAction, audit.Attempted, nil); auditErr != nil {
		http.Error(w, auditErr.Error(), http.StatusInternalServerError)
		return
	}

	response, err := api.getIdentity(ctx, r)

	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "getIdentity: error"), nil)
		api.auditor.Record(ctx, getIdentityAction, audit.Unsuccessful, nil)
		getIdentityResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, getIdentityAction, audit.Successful, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getIdentityResponse.writeEntity(ctx, w, response, http.StatusOK)
	log.InfoCtx(ctx, "getIdentity: get identity successful", log.Data{"ID": response.ID})
}

func (api *API) getIdentity(ctx context.Context, r *http.Request) (*schema.Identity, error) {

	tokenStr := r.Header.Get(tokenHeaderKey)
	if tokenStr == "" {
		log.ErrorCtx(ctx, ErrNoTokenProvided, nil)
		return nil, ErrNoTokenProvided
	}

	i, err := api.IdentityService.Get(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	return i, nil
}
