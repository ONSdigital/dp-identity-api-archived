package mongo

import (
	"github.com/ONSdigital/dp-identity-api/schema"
	"time"
)

func (m *Mongo) StoreToken(key string, i schema.Identity, ttl time.Duration) error {
	return nil
}

func (m *Mongo) GetToken(tokenStr string) (time.Duration, error) {
	return time.Second * 0, nil
}
