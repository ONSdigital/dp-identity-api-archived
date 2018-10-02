package mongo

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"time"
)

const (
	identityIDKey = "identity_id"
)

// StoreToken create a new token document in mongodb tokens collection. Any active token associated with the identity
// will be marked as deleted. Sets the last modified date on all documents updated.
func (m *Mongo) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error {
	logD := log.Data{identityIDKey: i.ID}
	log.InfoCtx(ctx, "tokenStore: storing identity token", logD)

	active, err := m.findActiveTokens(ctx, i)
	if err != nil {
		return errors.Wrap(err, "tokenStore findActiveTokens return an error")
	}

	err = m.deleteActiveTokens(ctx, i, active)
	if err != nil {
		return err
	}

	s := m.Session.Copy()
	defer s.Close()

	err = m.storeNewActiveToken(ctx, tkn)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) GetToken(ctx context.Context, tokenStr string) (time.Duration, error) {
	return time.Second * 0, nil
}

func (m *Mongo) findActiveTokens(ctx context.Context, i schema.Identity) ([]schema.Token, error) {
	log.InfoCtx(ctx, "tokenStore: querying for active tokens", log.Data{identityIDKey: i.ID})
	s := m.Session.Copy()
	defer s.Close()

	query := bson.M{"identity_id": i.ID, "deleted": false}

	var active []schema.Token
	err := s.DB(m.Database).C(m.TokenCollection).Find(query).All(&active)
	if err != nil {
		return nil, errors.Wrap(err, "tokenStore: query for active tokens returned an error")
	}
	return active, nil
}

func (m *Mongo) deleteActiveTokens(ctx context.Context, i schema.Identity, active []schema.Token) error {
	logD := log.Data{identityIDKey: i.ID}

	if len(active) == 0 {
		log.InfoCtx(ctx, "tokenStore: no currently active tokens to delete for identity", logD)
		return nil
	}

	log.InfoCtx(ctx, "tokenStore: deleting active token(s) for identity", logD)

	s := m.Session.Copy()
	defer s.Close()

	selector := bson.M{"identity_id": i.ID}
	update := bson.M{"$set": bson.M{"deleted": true, "last_modified": time.Now()}}

	info, err := s.DB(m.Database).C(m.TokenCollection).UpdateAll(selector, update)
	if err != nil {
		return errors.Wrap(err, "tokenStore: error deleting active token(s) for identity")
	}

	log.InfoCtx(ctx, "tokenStore: delete active tokens completed without error", log.Data{
		"identity_id": i.ID,
		"changeInfo": changeInfo{
			"matched": info.Matched,
			"updated": info.Updated,
		},
	})
	return nil
}

func (m *Mongo) storeNewActiveToken(ctx context.Context, tkn schema.Token) error {
	s := m.Session.Copy()
	defer s.Close()

	logD := log.Data{identityIDKey: tkn.IdentityID}
	log.InfoCtx(ctx, "tokenStore: storing new active identity token", logD)

	tkn.LastModified = time.Now()

	err := s.DB(m.Database).C(m.TokenCollection).Insert(tkn)
	if err != nil {
		return errors.Wrap(err, "tokenStore: error while storing new active identity token")
	}

	log.InfoCtx(ctx, "tokenStore: store new active identity token successful", logD)
	return nil
}
