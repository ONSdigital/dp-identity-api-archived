package api

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/audit"
	"net/http"
	"time"
)

//go:generate moq -out generate_mocks.go -pkg api . IdentityService

const (
	createIdentityAction = "createIdentity"
)

var (
	ErrFailedToReadRequestBody      = errors.New("error while attempting to read request body")
	ErrFailedToUnmarshalRequestBody = errors.New("error while attempting to unmarshal request body")

	//map identity errors to http status codes.
	errorStatusMapping = map[error]int{
		ErrFailedToUnmarshalRequestBody: http.StatusInternalServerError,
		ErrFailedToReadRequestBody:      http.StatusInternalServerError,
		identity.ErrInvalidArguments:    http.StatusInternalServerError,
		identity.ErrPersistence:         http.StatusInternalServerError,
	}
)

//API defines HTTP HandlerFunc's for the endpoints offered by the Identity API service.
type API struct {
	IdentityService    IdentityService
	healthCheckTimeout time.Duration
	auditor            audit.AuditorService
}

//IdentityService is a service for creating, updating and deleting Identities.
type IdentityService interface {
	Create(ctx context.Context, i *identity.Model) error
}
