package store

import "github.com/ONSdigital/dp-identity-api/mongo"

// DataStore provides a datastore.Storer struct to attach methods to
type DataStore struct {
	Backend mongo.Mongo
}

/*
Methods for querying mongo, add as needed

example .....
-----------------------

func (s *DataSTore) GetIdentity(id) (*models.identity, error) {

	// some code to query mongo and get id

return identity, nil
}

------------------------

accessed in the handlers via something like:

identity :=  api.storer.GetIdentity(id)

*/
