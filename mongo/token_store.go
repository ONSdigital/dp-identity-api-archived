package mongo

import (
	"context"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/globalsign/mgo/bson"
	"time"
)

func (m *Mongo) StoreToken(ctx context.Context, tkn schema.Token, i schema.Identity, ttl time.Duration) error {
	active, err := m.findActiveTokens(i)
	if err != nil {
		return nil
	}

	err = m.deleteActiveTokens(i, active)
	if err != nil {
		return err
	}

	s := m.Session.Copy()
	defer s.Close()

	err = m.storeNewActiveToken(tkn)
	if err != nil {
		return err
	}
}

func (m *Mongo) GetToken(ctx context.Context, tokenStr string) (time.Duration, error) {
	return time.Second * 0, nil
}

func (m *Mongo) findActiveTokens(i schema.Identity) ([]schema.Token, error) {
	s := m.Session.Copy()
	defer s.Close()

	query := bson.M{"identity_id": i.ID, "deleted": false}

	var active []schema.Token
	err := s.DB(m.Database).C(m.TokenCollection).Find(query).All(&active)
	if err != nil {
		return nil, err
	}
	return active, nil
}

func (m *Mongo) deleteActiveTokens(i schema.Identity, active []schema.Token) error {
	if len(active) == 0 {
		return nil
	}

	s := m.Session.Copy()
	defer s.Close()

	log.Info("deleting active tokens for identity", log.Data{"identity_id": i.ID})

	selector := bson.M{"identity_id": i.ID}
	update := bson.M{"$set": bson.M{"deleted": true, "last_modified": time.Now()}}

	info, err := s.DB(m.Database).C(m.TokenCollection).UpdateAll(selector, update)
	if err != nil {
		return err
	}

	log.Info("delete active tokens completed without error", log.Data{
		"identity_id": i.ID,
		"changeInfo": changeInfo{
			"matched": info.Matched,
			"updated": info.Updated,
		},
	})
	return nil
}

func (m *Mongo) storeNewActiveToken(tkn schema.Token) error {
	s := m.Session.Copy()
	defer s.Close()

	tkn.LastModified = time.Now()

	err := s.DB(m.Database).C(m.TokenCollection).Insert(tkn)
	if err != nil {
		return err
	}

	return nil
}
