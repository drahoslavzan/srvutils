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

func (m *LRUCache[T]) Get(key string) (value T, ok bool) {
	return m.cache.Get(key)
}

func (m *LRUCache[T]) Set(key string, value T) {
	m.cache.Add(key, value)
}

func (m *LRUCache[T]) Delete(key string) {
	m.cache.Remove(key)
}
