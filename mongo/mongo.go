package mongo

import (
	"github.com/globalsign/mgo"
	"time"

	"errors"
	"github.com/ONSdigital/dp-identity-api/models"
	"github.com/globalsign/mgo/bson"
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	Collection     string
	Database       string
	Session        *mgo.Session
	URI            string
	lastPingTime   time.Time
	lastPingResult error
}

// Init creates a new mgo.Session with a strong consistency and a write mode of "majortiy".
func (m *Mongo) Init() (session *mgo.Session, err error) {
	if session != nil {
		return nil, errors.New("session already exists")
	}

	if session, err = mgo.Dial(m.URI); err != nil {
		return nil, err
	}

	session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	session.SetMode(mgo.Strong, true)
	return session, nil
}

// GetIdentity gets an identity document
func (store *Mongo) GetIdentity(id string) (*models.Identity, error) {

	s := store.Session.Copy()
	defer s.Close()
	var identity models.Identity
	err := s.DB(store.Database).C("identities").Find(bson.M{"_id": id}).One(&identity)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("identity not found")
		}
		return nil, err
	}

	return &identity, nil
}

// CreateIdentity creates an identty document
func (store *Mongo) CreateIdentity(identity *models.Identity) error {

	s := store.Session.Copy()
	defer s.Close()

	err := s.DB(store.Database).C("identities").Insert(identity)
	if err == mgo.ErrNotFound {
		return errors.New("Failed to post identity to mongo")
	}

	if err != nil {
		return err
	}

	return nil
}
