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

func (api *API) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenReq, err := getNewTokenRequest(ctx, r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "authentication unsuccessful"), nil)
		newTokenResponse.writeError(ctx, w, err)
		return
	}

	p := common.Params{"email": tokenReq.Email}
	logD := log.Data{"email": tokenReq.Email}

	if auditErr := api.auditor.Record(ctx, createToken, audit.Attempted, p); auditErr != nil {
		newTokenResponse.writeError(ctx, w, auditErr)
		return
	}

	authToken, err := api.createToken(ctx, tokenReq)

	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "createToken: returned error"), logD)
		if auditErr := api.auditor.Record(ctx, createToken, audit.Unsuccessful, p); auditErr != nil {
			err = auditErr
		}
		newTokenResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, createToken, audit.Successful, p)
	if err != nil {
		newTokenResponse.writeError(ctx, w, ErrInternalServerError)
		return
	}

	log.InfoCtx(ctx, "createToken: request successful", logD)
	newTokenResponse.writeEntity(ctx, w, authToken, http.StatusOK)
}

func (api *API) createToken(ctx context.Context, tokenReq *NewTokenRequest) (*AuthToken, error) {
	logD := log.Data{"email": tokenReq.Email}

	i, err := api.IdentityService.VerifyPassword(ctx, tokenReq.Email, tokenReq.Password)
	if err != nil {
		log.ErrorCtx(ctx, errors.Wrap(err, "createToken: request unsuccessful"), logD)
		return nil, err
	}

	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	api.Cache.Set(token.String(), *i)

	log.InfoCtx(ctx, "createToken: user credential successfully verified", logD)
	return &AuthToken{Token: token.String()}, nil
}
