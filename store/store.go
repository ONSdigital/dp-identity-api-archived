package store

import "github.com/ONSdigital/dp-identity-api/models"

//go:generate moq -out storetest/generate_mocks.go -pkg storetest . Storer

type DataStore struct {
	Backend Storer
}

// Storer represents basic data access
type Storer interface {
	CreateIdentity(identity *models.Identity) error
}
