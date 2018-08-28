package store

import (
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/ONSdigital/dp-identity-api/mongo"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

type DataStore struct {
	Backend mongo.Mongo
}


func (store *DataStore) GetIdentity(id string) (*models.Identity, error) {

	s := store.Backend.Session.Copy()
	defer s.Close()
	var identity models.Identity
	err := s.DB(store.Backend.Database).C("identities").Find(bson.M{"_id": id}).One(&identity)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("identity not found")
		}
		return nil, err
	}

	return &identity, nil
}


func (store *DataStore) CreateIdentity(identity *models.Identity) error {

	s := store.Backend.Session.Copy()
	defer s.Close()

	err := s.DB(store.Backend.Database).C("identities").Insert(identity)
	if err == mgo.ErrNotFound {
		return errors.New("Failed to post identity to mongo")
	}

	if err != nil {
		return err
	}

	return nil
}
