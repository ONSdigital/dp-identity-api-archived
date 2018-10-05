package mongo

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/persistence"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"time"
)

const (
	identityIDKey = "identity_id"
)

// StoreToken store a new token document in mongodb tokens collection. Any active token associated with the identity
// will be marked as deleted. Sets the last modified date on all documents updated.
func (m *Mongo) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity) error {
	logD := log.Data{identityIDKey: i.ID}
	log.InfoCtx(ctx, "tokenStore: storing identity token", logD)

	_, err := m.deleteTokens(ctx, i.ID)
	if err != nil {
		return errors.Wrap(err, "error deleting tokens")
	}

	err = m.storeNewActiveToken(ctx, tkn)
	if err != nil {
		return errors.Wrap(err, "error storing new active token")
	}
	return nil
}

func (m *Mongo) GetIdentityByToken(ctx context.Context, token string) (*schema.Identity, *schema.Token, error) {
	s := m.Session.Copy()
	defer s.Close()

	t, err := m.getTokenByID(ctx, token)
	if err != nil && err == persistence.ErrNotFound { // token does not exit
		return nil, nil, err
	}

	// some other error querying fot token.
	if err != nil {
		return nil, nil, errors.Wrap(err, "error getting token by ID")
	}

	i, err := m.GetIdentityByID(ctx, t.IdentityID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error getting identity by ID")
	}
	return i, t, nil
}

// deleteTokens soft delete any active token associated with the provided identity ID. Sets token.deleted = true and
// updates token.last_modified to the current time. Returns the number of documents updated or any error encountered
// while executing.
func (m *Mongo) deleteTokens(ctx context.Context, identityID string) (int, error) {
	logD := log.Data{identityIDKey: identityID}
	log.InfoCtx(ctx, "tokenStore: deleting active token(s) for identity", logD)

	s := m.Session.Copy()
	defer s.Close()

	selector := bson.M{"identity_id": identityID, "deleted": false}
	update := bson.M{"$set": bson.M{"deleted": true, "last_modified": time.Now()}}

	info, err := s.DB(m.Database).C(m.TokenCollection).UpdateAll(selector, update)
	if err != nil {
		return 0, errors.Wrap(err, "tokenStore: error deleting active token(s) for identity")
	}

	log.InfoCtx(ctx, "tokenStore: delete active tokens completed without error", log.Data{
		"identity_id": identityID,
		"changeInfo": changeInfo{
			"matched": info.Matched,
			"updated": info.Updated,
		},
	})
	return info.Updated, nil
}

// storeNewActiveToken store the provided token in the Tokens collection. Token will become the active token for this
// identity.
func (m *Mongo) storeNewActiveToken(ctx context.Context, tkn schema.Token) error {
	s := m.Session.Copy()
	defer s.Close()

	logD := log.Data{identityIDKey: tkn.IdentityID}
	log.InfoCtx(ctx, "tokenStore: storing new active identity token", logD)

	tkn.LastModified = time.Now()
	tkn.Deleted = false // always set to false to ensure this is now the active token.

	err := s.DB(m.Database).C(m.TokenCollection).Insert(tkn)
	if err != nil {
		return errors.Wrap(err, "tokenStore: error while storing new active identity token")
	}

	log.InfoCtx(ctx, "tokenStore: store new active identity token successful", logD)
	return nil
}

func (m *Mongo) getTokenByID(ctx context.Context, tokenID string) (*schema.Token, error) {
	s := m.Session.Copy()
	defer s.Close()

	queryForToken := bson.M{"token_id": tokenID, "deleted": false}

	var t schema.Token
	if err := s.DB(m.Database).C(m.TokenCollection).Find(queryForToken).One(&t); err != nil {
		if err == mgo.ErrNotFound {
			log.InfoCtx(ctx, "active token for this values does not exist", nil)
			return nil, persistence.ErrNotFound
		}
		return nil, errors.Wrap(err, "error querying for active token")
	}
	return &t, nil
}

// activeTokens return a list of tokens associated with the provided identity with "deleted = false".
func (m *Mongo) getActiveTokensByIdentity(ctx context.Context, identityID string) ([]schema.Token, error) {
	log.InfoCtx(ctx, "tokenStore: querying for active tokens", log.Data{identityIDKey: identityID})
	s := m.Session.Copy()
	defer s.Close()

	query := bson.M{"identity_id": identityID, "deleted": false}

	var active []schema.Token
	err := s.DB(m.Database).C(m.TokenCollection).Find(query).All(&active)
	if err != nil {
		return nil, errors.Wrap(err, "tokenStore: query for active tokens returned an error")
	}

	return active, nil
}
