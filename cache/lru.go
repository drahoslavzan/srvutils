package cache

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"
)

type (
	LRUCache[T any] struct {
		cache *lru.Cache[string, T]
	}
)

func NewLRUCache[T any](sz int) *LRUCache[T] {
	cache, err := lru.New[string, T](sz)
	if err != nil {
		zap.L().Panic("cannot initialize lru cache", zap.Error(err))
	}

	return &LRUCache[T]{
		cache: cache,
	}
}

func (m *LRUCache[T]) Get(key string) *T {
	if v, ok := m.cache.Get(key); ok {
		return &v
	}

	return nil
}

func (m *LRUCache[T]) Set(key string, value T) {
	m.cache.Add(key, value)
}
