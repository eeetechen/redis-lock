package redis_lock

import (
	"github.com/go-redis/redis/v8"
	"redis-lock/distributed_lock"
)

type Client struct {
	Rdb  *redis.Client
	lock *distributed_lock.RedisLock
}

type Option struct {
	Addr     string
	Password string
	DB       int
}

func NewClient(option Option) Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: option.Password, // no password set
		DB:       option.DB,       // use default DB
	})

	return Client{
		Rdb: rdb,
	}

}

func (c *Client) GetRedisReadLock(name string) {
	c.lock = distributed_lock.NewRedisLock(name, distributed_lock.RedisReadLocked)
}

func (c *Client) GetRedisWriteLock(name string) {
	c.lock = distributed_lock.NewRedisLock(name, distributed_lock.RedisWriteLocked)
}

func (c *Client) ClearLock() {
	c.lock = nil
}
