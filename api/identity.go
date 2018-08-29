package api

import (
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/go-ns/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)


func (api *IdentityAPI) CreateIdentity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to read request body"), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var identity *models.Identity
	err = json.Unmarshal(body, &identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to unmarshall request body"), nil)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.dataStore.Backend.CreateIdentity(identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to write data to mongo"), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
