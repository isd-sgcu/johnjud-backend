package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	SetValue(key string, value interface{}, ttl int) error
	GetValue(key string, value interface{}) error
	DeleteValue(key string) error
}

type repositoryImpl struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) Repository {
	return &repositoryImpl{client: client}
}

func (r *repositoryImpl) SetValue(key string, value interface{}, ttl int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, v, time.Duration(ttl)*time.Second).Err()
}

func (r *repositoryImpl) GetValue(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	v, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(v), value)
}

func (r *repositoryImpl) DeleteValue(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Del(ctx, key).Err()
}
