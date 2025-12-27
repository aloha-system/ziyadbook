package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"ziyadbook/internal/config"
)

func New(cfg config.Config) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.Ping(ctx).Err(); err != nil {
		_ = c.Close()
		return nil, err
	}

	return c, nil
}
