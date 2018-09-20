package cache

import "github.com/ONSdigital/go-ns/log"

// NOPCache is a no op implementation of a cache.
type NOPCache struct {}


func (c *NOPCache) Set(key string, i interface{}) error {
	log.Info("nopcache: set", log.Data{"key": key})
	return nil
}

func (c *NOPCache) Get(key string, i interface{}) error {
	log.Info("nopcache: get", log.Data{"key": key})
	return nil
}

func (c *NOPCache) Delete(key string) error {
	log.Info("nopcache: delete", log.Data{"key": key})
	return nil
}
