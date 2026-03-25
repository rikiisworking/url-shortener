package storage

import (
	"context"
	"time"

	"gitgub.com/rikiisworking/url-shortener/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func NewRedisRepo(cfg config.Config) *RedisRepo {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	return &RedisRepo{Client: client}
}

func (r *RedisRepo) Set(shortCode, originalURL string, expiration time.Duration) error {
	return r.Client.Set(context.Background(), shortCode, originalURL, expiration).Err()
}

func (r *RedisRepo) Get(shortCode string) (string, error) {
	return r.Client.Get(context.Background(), shortCode).Result()
}

func (r *RedisRepo) IncrementClick(shortCode string) error {
	return r.Client.Incr(context.Background(), shortCode+":clicks").Err()
}
