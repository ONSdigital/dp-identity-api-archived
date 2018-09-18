package cache

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrTokenNotFound = errors.New("could not find token in cache")
)

func New(addr string) *IdentityCache {
	pool := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			//return redis.Dial("tcp", ":6379")
			return redis.Dial("tcp", addr)
		},
	}

	return &IdentityCache{
		ttl:  30,
		pool: pool,
	}

}

type IdentityCache struct {
	pool redis.Pool
	ttl  int
}

func (c *IdentityCache) Set(token string, i identity.Model) error {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HMSET", redis.Args{token}.AddFlat(i)...)
	conn.Send("EXPIRE", token, c.ttl)
	r, err := conn.Do("EXEC")
	if err != nil {
		return err
	}
	fmt.Println(r)

	return nil
}

func (c *IdentityCache) Get(token string) (*identity.Model, error) {
	conn := c.pool.Get()
	defer conn.Close()

	// Create a transaction and queue the commands to update TTL and get the identity object.
	conn.Send("MULTI")
	conn.Send("EXPIRE", token, c.ttl)
	conn.Send("HGETALL", token)

	// execute the transaction
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	// get the EXPIRE response
	expire, err := redis.Int64(values[0], nil)
	if err != nil {
		return nil, err
	}

	// EXPIRE returns 1 if the TTL was set, 0 if the key did not exist.
	if expire == 0 {
		log.Info("redis key does not exist", nil)
		return nil, ErrTokenNotFound
	}

	// get the response to HGETALL i.e. the identity object.
	hgetall, err := redis.Values(values[1], nil)
	if err != nil {
		return nil, err
	}

	// Convert the response bytes into a struct
	var i identity.Model
	err = redis.ScanStruct(hgetall, &i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (c *IdentityCache) Close(ctx context.Context) error {
	log.Info("closing IdentityCache", nil)
	return c.pool.Close()
}
