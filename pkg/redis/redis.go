package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

const (
	maxRetries  = 5
	dialTimeout = 5 * time.Second
)

func NewRedisClient(cfg *Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:        cfg.Addr,
		Password:    cfg.Password,
		DB:          cfg.DB,
		MaxRetries:  maxRetries,
		DialTimeout: dialTimeout,
	})

	return client
}
