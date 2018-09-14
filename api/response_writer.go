package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("internal server error")

	createIdentityResponse = JSONResponseWriter{
		ErrFailedToUnmarshalRequestBody: http.StatusBadRequest,
		ErrFailedToReadRequestBody:      http.StatusBadRequest,
		ErrRequestBodyNil:               http.StatusBadRequest,
		identity.ErrInvalidArguments:    http.StatusInternalServerError,
		identity.ErrPersistence:         http.StatusInternalServerError,
		identity.ErrNameValidation:      http.StatusBadRequest,
		identity.ErrEmailValidation:     http.StatusBadRequest,
		identity.ErrPasswordValidation:  http.StatusBadRequest,
		identity.ErrIdentityNil:         http.StatusBadRequest,
	}

	getIdentityResponse = JSONResponseWriter{
		ErrNoTokenProvided:	http.StatusNotFound,
	}

	newTokenResponse = JSONResponseWriter{
		ErrRequestBodyNil:              http.StatusBadRequest,
		ErrAuthRequestNil:              http.StatusBadRequest,
		ErrAuthRequestIDNil:            http.StatusBadRequest,
		identity.ErrAuthenticateFailed: http.StatusForbidden,
		identity.ErrUserNotFound:       http.StatusNotFound,
	}
)

type JSONResponseWriter map[error]int

func (e JSONResponseWriter) writeEntity(ctx context.Context, w http.ResponseWriter, i interface{}, status int) {
	b, err := json.Marshal(i)
	if err != nil {
		log.ErrorCtx(ctx, errors.New("failed to marshal object to JSON"), log.Data{"object": i})
		e.writeError(ctx, w, ErrInternalServerError)
		return
	}
	w.Header().Set(headerContentType, mimeTypeJSON)
	w.WriteHeader(status)
	w.Write(b) // TODO handle error
}

func (e JSONResponseWriter) writeError(ctx context.Context, w http.ResponseWriter, err error) {
	status := e.resolveError(err)

	if status == http.StatusInternalServerError {
		err = ErrInternalServerError
	}

	log.ErrorCtx(ctx, errors.New("writing error response"), log.Data{"status": status})
	http.Error(w, err.Error(), status)
}

func (e JSONResponseWriter) resolveError(err error) int {
	status := http.StatusInternalServerError

	if val, ok := e[err]; ok {
		status = val
	}
	return status
}
