package mongo

import (
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/globalsign/mgo"
	"time"

	"errors"
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


func (m *Mongo) Create(identity *identity.Model) error {
	s := m.Session.Copy()
	defer s.Close()

	// NOTE - Upsert may be more appropriate than Insert. Consider "already exists" scenarios?
	err := s.DB(m.Database).C("identities").Insert(identity)
	if err == mgo.ErrNotFound {
		return errors.New("failed to post new identity document to mongo")
	}

	if err != nil {
		return err
	}

	return nil
}
