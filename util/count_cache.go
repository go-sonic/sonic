package util

import (
	"sync"
	"time"
)

type CounterCache[K comparable] struct {
	rwMu            sync.RWMutex
	countCache      map[K]int64
	batchIncr       func(countCache map[K]int64)
	singleIncr      func(key K, count int64)
	refreshDuration time.Duration
}

func NewCounterCache[K comparable](refreshDuration time.Duration, batchIncr func(map[K]int64), singleIncr func(K, int64)) *CounterCache[K] {
	c := &CounterCache[K]{
		countCache:      make(map[K]int64),
		batchIncr:       batchIncr,
		singleIncr:      singleIncr,
		refreshDuration: refreshDuration,
	}
	go c.startFlushTicker()
	return c
}

func (c *CounterCache[K]) IncrBy(key K, value int64) int64 {
	val := c.incrCacheBy(key, value)
	return val
}

func (c *CounterCache[K]) incrCacheBy(key K, value int64) int64 {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	count := c.countCache[key]
	count += value
	c.countCache[key] = count
	return count
}

func (c *CounterCache[K]) Get(key K) int64 {
	cacheVal := c.get(key)
	return cacheVal
}

func (c *CounterCache[K]) get(key K) int64 {
	c.rwMu.RLock()
	defer c.rwMu.RUnlock()
	return c.countCache[key]
}

func (c *CounterCache[K]) startFlushTicker() {
	ticker := time.NewTicker(c.refreshDuration)
	defer ticker.Stop()
	for range ticker.C {
		c.flush()
	}
}

func (c *CounterCache[K]) flush() {
	var oldCountCache map[K]int64
	c.rwMu.Lock()
	oldCountCache = c.countCache
	c.countCache = make(map[K]int64)
	c.rwMu.Unlock()
	if c.batchIncr != nil {
		c.batchIncr(oldCountCache)
		return
	}
	for key, value := range oldCountCache {
		c.singleIncr(key, value)
	}
}
