package mongo

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"time"
)

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

	err = s.DB(m.Database).C(m.IdentityCollection).Insert(identity)
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

	count, err := s.DB(m.Database).C(m.IdentityCollection).Find(query).Count()
	if err != nil {
		return false, errors.Wrap(err, "error executing count active identities query")
	}

	return count == 0, nil
}

func (m *Mongo) GetIdentity(email string) (schema.Identity, error) {
	s := m.Session.Copy()
	defer s.Close()

	query := bson.M{"email": email, "deleted": false}

	count, err := s.DB(m.Database).C(m.IdentityCollection).Find(query).Count()
	if err != nil {
		return schema.NilIdentity, err
	}

	if count == 0 {
		return schema.NilIdentity, persistence.ErrNotFound
	}

	var i schema.Identity
	if err := s.DB(m.Database).C(m.IdentityCollection).Find(query).One(&i); err != nil {
		return schema.NilIdentity, err
	}
	return i, nil
}

func (m *Mongo) GetIdentityByID(ctx context.Context, id string) (*schema.Identity, error) {
	s := m.Session.Copy()
	defer s.Close()

	var i schema.Identity

	query := bson.M{"id": id, "deleted": false}

	if err := s.DB(m.Database).C(m.IdentityCollection).Find(query).One(&i); err != nil {
		if err == mgo.ErrNotFound {
			err = persistence.ErrNotFound
		}
		return nil, err
	}
	return &i, nil
}
