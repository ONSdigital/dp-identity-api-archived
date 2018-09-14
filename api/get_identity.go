package api

import (
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"net/http"
	"context"
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

	response, err := api.getIdentity(ctx)

	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "getIdentity: error"), nil)
		api.auditor.Record(ctx, getIdentityAction, audit.Unsuccessful, nil)
		getIdentityResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, getIdentityAction, audit.Successful, common.Params{"id": response.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getIdentityResponse.writeEntity(ctx, w, response, http.StatusOK)
	log.InfoCtx(ctx, "createIdentity: get identity successful", log.Data{"id": response.ID})
}

func (api *API) getIdentity(ctx context.Context) (*identity.Model, error) {

	tokenStr, ok := ctx.Value("token").(string)
	if !ok {
		return nil, identity.ErrNoTokenProvided
	}

	i, err := api.IdentityService.Get(tokenStr)
	if err != nil {
		return nil, err
	}

	return i, nil
}
