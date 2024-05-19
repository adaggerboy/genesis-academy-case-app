package openexchangeapi

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(key K, reloadFunc func(K) (V, error)) (V, error)
}

type cacheDeadlineValueHolder[V any] struct {
	value    V
	deadline time.Time
}

type cacheImpl[K comparable, V any] struct {
	Cache[K, V]
	sync.RWMutex
	internalMap map[K]cacheDeadlineValueHolder[V]
	defaultTTL  time.Duration
}

func (c *cacheImpl[K, V]) reloadGet(key K, reloadFunc func(K) (V, error)) (value V, err error) {
	v, err := reloadFunc(key)
	if err != nil {
		return
	}
	holder := cacheDeadlineValueHolder[V]{
		deadline: time.Now().Add(c.defaultTTL),
		value:    v,
	}
	c.Lock()
	c.internalMap[key] = holder
	c.Unlock()
	value = holder.value
	return
}

func (c *cacheImpl[K, V]) Get(key K, reloadFunc func(K) (V, error)) (value V, err error) {
	c.RLock()
	holder, ok := c.internalMap[key]
	c.RUnlock()
	if !ok {
		value, err = c.reloadGet(key, reloadFunc)
	} else if holder.deadline.Before(time.Now()) {
		value, err = c.reloadGet(key, reloadFunc)
	} else {
		value = holder.value
		err = nil
	}
	return
}

func NewCache[K comparable, V any](duration time.Duration) Cache[K, V] {
	return &cacheImpl[K, V]{
		internalMap: map[K]cacheDeadlineValueHolder[V]{},
		defaultTTL:  duration,
	}
}
