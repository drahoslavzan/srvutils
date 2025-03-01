package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	ValueFactory[T any] interface {
		FromString(string) (T, error)
		ToString(T) string
	}

	KVCache[T any] struct {
		client *redis.Client
		ttl    time.Duration
		vf     ValueFactory[T]
		logger *zap.Logger
	}
)

const NoTTL = time.Duration(0)

func NewRedisCache[T any](conn string, ttl time.Duration, vf ValueFactory[T], logger *zap.Logger) *KVCache[T] {
	opts, err := redis.ParseURL(conn)
	if err != nil {
		logger.Panic("redis parse", zap.Error(err))
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

func (m *KVCache[T]) Get(key string) (value T, ok bool) {
	v, err := m.client.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil {
			m.logger.Error("redis get", zap.Error(err))
		}

		return
	}

	ret, err := m.vf.FromString(v)
	if err != nil {
		m.logger.Error("value factory from string", zap.Error(err))
		return
	}

	return ret, true
}

func (m *KVCache[T]) Set(key string, value T) {
	err := m.client.Set(context.Background(), key, m.vf.ToString(value), m.ttl).Err()
	if err != nil {
		m.logger.Error("redis set", zap.Error(err))
		return
	}
}

func (m *KVCache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	err := m.client.Set(context.Background(), key, value, ttl).Err()
	if err != nil {
		m.logger.Error("redis set", zap.Error(err))
		return
	}
}

func (m *KVCache[T]) Delete(key string) {
	err := m.client.Del(context.Background(), key).Err()
	if err != nil {
		m.logger.Error("redis delete", zap.Error(err))
		return
	}
}
