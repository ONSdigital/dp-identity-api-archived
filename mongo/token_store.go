package mongo

import (
	"time"
	"github.com/ONSdigital/dp-identity-api/schema"
)

func (m *Mongo) StoreToken(key string, i schema.Identity, ttl time.Duration) error {
	return nil
}

func (m *Mongo) GetToken(token string) (time.Duration, error) {
	return time.Second * 0, nil
}
