package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (api *API) ping(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	_, err := api.Cache.Get(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
