package api

import (
	"context"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"net/http"
)

func (api *API) AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authReq, err := getAuthenticateRequest(r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authentication unsuccessful"), nil)
		authenticateResponse.writeError(ctx, w, err)
		return
	}

	p := common.Params{"id": authReq.ID}
	logD := log.Data{"id": authReq.ID}

	if auditErr := api.auditor.Record(ctx, authenticateAction, audit.Attempted, p); auditErr != nil {
		http.Error(w, auditErr.Error(), http.StatusInternalServerError)
		return
	}

	authToken, err := api.authenticate(ctx, authReq)

	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authenticateResponse: returned error"), logD)
		api.auditor.Record(ctx, authenticateAction, audit.Unsuccessful, p)
		authenticateResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, authenticateAction, audit.Successful, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.InfoCtx(ctx, "authenticateResponse: request successful", logD)
	authenticateResponse.writeEntity(ctx, w, authToken, http.StatusOK)
}

func (api *API) authenticate(ctx context.Context, authReq *AuthenticateRequest) (*AuthToken, error) {
	logD := log.Data{"id": authReq.ID}

	err := api.IdentityService.Authenticate(ctx, authReq.ID, authReq.Password)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authenticate: request unsuccessful"), logD)
		return nil, err
	}

	log.ErrorCtx(ctx, errors.Wrap(err, "authenticate: request successful"), logD)
	return &AuthToken{Token: "1234567890"}, nil
}
