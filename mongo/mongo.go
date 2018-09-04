package mongo

import (
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/globalsign/mgo"
	"github.com/satori/go.uuid"
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

func New(cfg config.MongoConfig) (*Mongo, error) {
	mongodb := &Mongo{
		Collection: cfg.Collection,
		Database:   cfg.Database,
		URI:        cfg.BindAddr,
	}

	session, err := mongodb.createSession()
	if err != nil {
		return nil, err
	}

	mongodb.Session = session
	return mongodb, nil
}

// createSession creates a new mgo.Session with a strong consistency and a write mode of "majortiy".
func (m *Mongo) createSession() (session *mgo.Session, err error) {
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

func (m *Mongo) Create(identity *identity.Model) (string, error) {
	s := m.Session.Copy()
	defer s.Close()

	// NOTE - Upsert may be more appropriate than Insert. Consider "already exists" scenarios?
	id := uuid.NewV4()
	identity.ID = id.String()

	err := s.DB(m.Database).C("identities").Insert(identity)
	if err == mgo.ErrNotFound {
		return "", errors.New("failed to post new identity document to mongo")
	}

	if err != nil {
		return "", err
	}

	return id.String(), nil
}
