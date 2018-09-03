package api

import (
	"context"
	"net/http"
)

//go:generate moq -out generate_mocks.go -pkg api . IdentityService

//IdentityService is a service for creating, updating and deleting Identities.
type IdentityService interface {
	Create(ctx context.Context, r *http.Request) error
}