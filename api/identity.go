package api

import (
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"context"
)

const (
	createIdentityAction = "createIdentity"
)

var (
	ErrFailedToReadRequestBody = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to read request body",
	}

	ErrFailedToMarshalRequestBody = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to unmarshall request body",
	}

	ErrFailedToWriteToMongo = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to write data to mongo",
	}
)

func (api *IdentityAPI) createIdentity(ctx context.Context, r *http.Request) *apiError {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "createIdentity error failed to read request body"), nil)
		return ErrFailedToReadRequestBody
	}

	var identity *models.Identity
	err = json.Unmarshal(body, &identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "createIdentity: failed to unmarshall request body"), nil)
		return ErrFailedToMarshalRequestBody
	}

	err = api.dataStore.Backend.CreateIdentity(identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "createIdentity: failed to write data to mongo"), nil)
		return ErrFailedToWriteToMongo
	}
	return nil
}
