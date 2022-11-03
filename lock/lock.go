package lock

import (
	"math/rand"
	"sync"
	"time"
)

type RedisLock struct {
	name      string
	timestamp int64
	mu        sync.Mutex
	status    uint8
	randVal   int
}

func NewRedisLock(name string, status uint8) *RedisLock {
	return &RedisLock{
		name:      name,
		timestamp: time.Now().UnixNano(),
		mu:        sync.Mutex{},
		status:    status,
		randVal:   rand.Int(),
	}
}

func (r *RedisLock) TryLock() bool {
	return r.mu.TryLock()
}

func (r *RedisLock) Lock() {
	r.mu.Lock()
}

func (r *RedisLock) Unlock() {
	r.mu.Unlock()
}
