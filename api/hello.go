package api

import (
	"github.com/ONSdigital/go-ns/audit"
	"github.com/gorilla/mux"
	"net/http"
)

// simple endpoint for checking functionality - delete as soon as appropriate
func (api *IdentityAPI) hello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	name := vars["name"]

	// See that the auditor works
	auditErr := api.auditor.Record(ctx, "hello action", audit.Attempted, nil)
	if auditErr != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello " + name))
}
