package provider

import (
	"github.com/redis/go-redis/v9"
	"gits/internal/container"
)

type Cache interface {
	GetConnection() *redis.Client
	CloseConnection() error
}

type cache struct {
	container container.Container
	client    *redis.Client
}

func NewCache(container container.Container) Cache {
	config := container.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Cache.Address(),
		Password: config.Cache.Password,
	})
	return &cache{
		container: container,
		client:    client,
	}
}

func (c *cache) GetConnection() *redis.Client {
	return c.client
}

func (c *cache) CloseConnection() error {
	return c.client.Close()
}
