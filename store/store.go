package store

import "github.com/ONSdigital/dp-identity-api/models"

type DataStore struct {
	Backend Storer
}

// Storer represents basic data access
type Storer interface {
	GetIdentity(id string) (*models.Identity, error)
	CreateIdentity(identity *models.Identity) error
}
