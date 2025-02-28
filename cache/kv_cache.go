package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	ValueFactory[T any] interface {
		FromString(string) (*T, error)
	}

	KVCache[T any] struct {
		client *redis.Client
		ttl    time.Duration
		vf     ValueFactory[T]
		logger *zap.Logger
	}
)

func NewRedisCache[T any](conn string, ttl time.Duration, vf ValueFactory[T], logger *zap.Logger) *KVCache[T] {
	opts, err := redis.ParseURL(conn)
	if err != nil {
		panic(err)
	}

	cln := redis.NewClient(opts)

	return &KVCache[T]{
		client: cln,
		ttl:    ttl,
		vf:     vf,
		logger: logger,
	}
}

func (m *KVCache[T]) Ping() error {
	return m.client.Ping(context.Background()).Err()
}

func (m *KVCache[T]) Get(key string) *T {
	v, err := m.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		m.logger.Error("redis get", zap.Error(err))
		return nil
	}

	ret, err := m.vf.FromString(v)
	if err != nil {
		m.logger.Error("value factory from string", zap.Error(err))
		return nil
	}

	return ret
}

func (m *KVCache[T]) Set(key string, value T) {
	err := m.client.Set(context.Background(), key, value, m.ttl).Err()
	if err != nil {
		m.logger.Error("redis set", zap.Error(err))
		return
	}
}
