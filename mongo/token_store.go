package mongo

import "time"

func (m *Mongo) StoreToken(token string, identityID string) (time.Duration, error) {
	return time.Second * 0, nil
}

func (m *Mongo) GetToken(token string, identityID string) (time.Duration, error) {
	return time.Second * 0, nil
}
