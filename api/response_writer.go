package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("internal server error")
)

func writeJSONBody(ctx context.Context, w http.ResponseWriter, i interface{}, status int) {
	b, err := json.Marshal(i)
	if err != nil {
		log.ErrorCtx(ctx, errors.New("failed to marshal object to JSON"), log.Data{"object": i})
		writeErrorResponse(ctx, ErrInternalServerError, w)
		return
	}
	w.Header().Set(headerContentType, mimeTypeJSON)
	w.WriteHeader(status)
	w.Write(b)
}

//writeErrorResponse writes a HTTP error back to the response writer. If the err can be cast to apiError then the values
// of err.GetMessage() and err.GetStatus() will be used to set the response body and status code respectively otherwise
// a default 500 status is used with err.Error() for the response body.
func writeErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError

	if val, ok := errorStatusMapping[err]; ok {
		status = val
	}

	log.ErrorCtx(ctx, errors.New("writing error response"), log.Data{"status": status})
	http.Error(w, err.Error(), status)
}
