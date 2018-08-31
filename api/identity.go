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
	ErrInvalidArguments = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting create new identity",
	}

	ErrFailedToReadRequestBody = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to read request body",
	}

	ErrFailedToUnmarshalRequestBody = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to unmarshal request body",
	}

	ErrFailedToWriteToMongo = &apiError{
		status:  http.StatusInternalServerError,
		message: "error while attempting to write data to mongo",
	}
)

//createIdentity contains the business logic for creating a new Identity.
func (api *IdentityAPI) createIdentity(ctx context.Context, r *http.Request) *apiError {
	if ctx == nil {
		log.Error(errors.New("createIdentity: failed mandatory context parameter was nil"), nil)
		return ErrInvalidArguments
	}

	if r == nil {
		log.Error(errors.New("createIdentity: failed mandatory request parameter was nil"), nil)
		return ErrInvalidArguments
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "createIdentity: failed to read request body"), nil)
		return ErrFailedToReadRequestBody
	}

	defer r.Body.Close()

	var identity *models.Identity
	err = json.Unmarshal(body, &identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "createIdentity: failed to unmarshal request body"), nil)
		return ErrFailedToUnmarshalRequestBody
	}

	err = api.dataStore.Backend.CreateIdentity(identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "createIdentity: failed to write data to mongo"), nil)
		return ErrFailedToWriteToMongo
	}
	return nil
}
