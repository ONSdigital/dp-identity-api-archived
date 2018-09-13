package api

import (
	"net/http"
	"github.com/pkg/errors"
	"context"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/common"
)

// TODO - meaningful documentation
//
//
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
		createIdentityResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, getIdentityAction, audit.Successful, common.Params{"id": response.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createIdentityResponse.writeEntity(ctx, w, response, http.StatusOK)
	log.InfoCtx(ctx, "createIdentity: get identity successful", log.Data{"id": response.ID})
}


func (api *API) getIdentity(ctx context.Context, r *http.Request) (*identity.Model, error) {

	i, err := api.IdentityService.Get(ctx)
	if err != nil {
		return nil, err
	}

	return i, nil
}