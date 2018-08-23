package api

import (
	"encoding/json"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/go-ns/audit"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// TODO - added to sanity check - remove/change/purge as needed
func (api *IdentityAPI) GetIdentityByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	identityID := vars["id"]

	logData := log.Data{"identity_id": identityID}
	auditParams := common.Params{"identity_id": identityID}

	if err := api.auditor.Record(ctx, "get Identity", audit.Attempted, auditParams); err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "request unsuccessful"), logData)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	identity, err := api.dataStore.GetIdentityByID(identityID)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to get identity"), logData)
		err := api.auditor.Record(ctx, "get Identity", audit.Unsuccessful, auditParams)
		if err != nil {
			log.ErrorCtx(ctx, errors.WithMessage(err, "request unsuccessful"), logData)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := api.auditor.Record(ctx, "get Identity", audit.Successful, auditParams); err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "request unsuccessful"), logData)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to marshal identity into bytes"), logData)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(b))
}

// TODO - added to sanity check - remove/change/purge as needed
func (api *IdentityAPI) PostIdentity(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = api.dataStore.PostIdentity(identity)
	if err != nil {
		log.ErrorCtx(ctx, errors.WithMessage(err, "failed to write data to mongo"), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
