package mongo

import (
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/globalsign/mgo"
	"time"

	"github.com/pkg/errors"
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
