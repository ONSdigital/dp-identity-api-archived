package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

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
		log.ErrorCtx(ctx, errors.Wrap(err, "createIdentity: error"), nil)
		api.auditor.Record(ctx, createIdentityAction, audit.Unsuccessful, nil)
		createIdentityResponse.writeError(ctx, w, err)
		return
	}

	err = api.auditor.Record(ctx, createIdentityAction, audit.Successful, common.Params{"id": response.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createIdentityResponse.writeEntity(ctx, w, response, http.StatusCreated)
	log.InfoCtx(ctx, "createIdentity: identity created successfully", log.Data{"id": response.ID})
}

func (api *API) createIdentity(ctx context.Context, r *http.Request) (*IdentityCreated, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, ErrFailedToReadRequestBody
	}
	defer r.Body.Close()

	if len(body) == 0 {
		return nil, ErrRequestBodyNil
	}

	var i schema.Identity
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
