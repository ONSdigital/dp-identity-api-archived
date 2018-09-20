package mongo

import (
	"github.com/ONSdigital/dp-identity-api/config"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
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

func (m *Mongo) SaveIdentity(identity schema.Identity) (string, error) {
	s := m.Session.Copy()
	defer s.Close()

	available, err := m.identityAvailable(s, identity.Email)
	if err != nil {
		return "", err
	}

	if !available {
		return "", persistence.ErrNonUnique
	}

	// NOTE - Upsert may be more appropriate than Insert. Consider "already exists" scenarios?
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(err, "error generating uuid")
	}

	identity.ID = id.String()
	identity.CreatedDate = time.Now()

	err = s.DB(m.Database).C(m.Collection).Insert(identity)
	if err == mgo.ErrNotFound {
		return "", errors.New("failed to post new identity document to mongo")
	}

	if err != nil {
		return "", err
	}

	return identity.ID, nil
}

func (m *Mongo) identityAvailable(s *mgo.Session, email string) (bool, error) {
	query := bson.M{"email": email, "deleted": false}

	count, err := s.DB(m.Database).C(m.Collection).Find(query).Count()
	if err != nil {
		return false, errors.Wrap(err, "error executing count active identities query")
	}

	return count == 0, nil
}

func (m *Mongo) GetIdentity(email string) (schema.Identity, error) {
	s := m.Session.Copy()
	defer s.Close()

	query := bson.M{"email": email, "deleted": false}

	count, err := s.DB(m.Database).C(m.Collection).Find(query).Count()
	if err != nil {
		return schema.NilIdentity, err
	}

	if count == 0 {
		return schema.NilIdentity, persistence.ErrNotFound
	}

	var i schema.Identity
	if err := s.DB(m.Database).C(m.Collection).Find(query).One(&i); err != nil {
		return schema.NilIdentity, err
	}
	return i, nil
}
