package api

import (
	"context"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"net/http"
)

var (
	ErrAuthRequestNil   = errors.New("authentication request invalid")
	ErrAuthRequestIDNil = errors.New("authentication request invalid: id required but was empty")
)

func (api *API) AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authReq, err := getAuthenticateRequest(ctx, r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authentication unsuccessful"), nil)
		authenticateResponse.writeError(ctx, w, err)
		return
	}

	p := common.Params{"id": authReq.ID}
	logD := log.Data{"id": authReq.ID}

	if auditErr := api.auditor.Record(ctx, authenticateAction, audit.Attempted, p); auditErr != nil {
		authenticateResponse.writeError(ctx, w, auditErr)
		return
	}

	authToken, err := api.authenticate(ctx, authReq)

	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authenticate: returned error"), logD)
		if auditErr := api.auditor.Record(ctx, authenticateAction, audit.Unsuccessful, p); auditErr != nil {
			err = auditErr
		}
		authenticateResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, authenticateAction, audit.Successful, p)
	if err != nil {
		authenticateResponse.writeError(ctx, w, ErrInternalServerError)
		return
	}

	log.InfoCtx(ctx, "authenticate: request successful", logD)
	authenticateResponse.writeEntity(ctx, w, authToken, http.StatusOK)
}

func (api *API) authenticate(ctx context.Context, authReq *AuthenticateRequest) (*AuthToken, error) {
	logD := log.Data{"id": authReq.ID}

	err := api.IdentityService.Authenticate(ctx, authReq.ID, authReq.Password)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authenticate: request unsuccessful"), logD)
		return nil, err
	}

	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	log.InfoCtx(ctx, "authenticate: user credential successfully verified", logD)
	return &AuthToken{Token: token.String()}, nil
}
