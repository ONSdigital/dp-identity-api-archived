package store

import "github.com/ONSdigital/dp-identity-api/models"

type DataStore struct {
	Backend Storer
}

// Storer represents basic data access
type Storer interface {
	CreateIdentity(identity *models.Identity) error
}
