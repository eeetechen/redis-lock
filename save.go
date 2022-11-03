package redis_lock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"redis-lock/lock"
	"time"

	"go.uber.org/zap"
)

// 写入
func WriteRedis(client Client, name, key, value string) error {

	//本地超时
	var ctx, cancel = context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()

	client.GetRedisWriteLock(name)
	defer client.ClearLock()
	client.lock.Lock()
	defer client.lock.Unlock()

	var pif = func(pipe redis.Pipeliner) error {
		cmd := pipe.Set(ctx, key, value, redis.KeepTTL)
		if cmd.Err() != nil {
			return cmd.Err()
		}
		return nil
	}
	var txf = func(tx *redis.Tx) error {
		cmds, err := tx.TxPipelined(ctx, pif)
		if err != nil {
			return err
		}

		for _, cmd := range cmds {
			cmd.
		}

		return nil
	}

	for i := 0; i < 10; i++ {
		err := client.Rdb.Watch(ctx, txf)
		if err == nil {
			// Success.
			return
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}

	cmd := client.Rdb.SetNX(ctx, name, lock.GenerateRedisLockVal(client.lock), time.Duration(15*time.Second))
	if cmd.Err() != nil {
		zap.S().Error(cmd.Err())
		return cmd.Err()
	}

	return nil
}

// 判断写锁的读取
func INDReadRedis(client Client, name, key string) (string, error) {
	//本地超时
	var ctx, cancel = context.WithTimeout(context.Background(), time.Duration(lock.RedisTTL)*time.Second)
	defer cancel()

	client.GetRedisReadLock(name)
	defer client.ClearLock()
	client.lock.Lock()
	defer client.lock.Unlock()

	retryCount := 0

Retry:
	cmd := client.Rdb.Get(ctx, name)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		zap.S().Error("Redis link Err. err:= %w", cmd.Err())
		return "", cmd.Err()
	}
	if cmd.Err() == nil {
		strVal := cmd.Val()
		strStatus, _, _, err := lock.ParseRedisLockVal(strVal)
		if err != nil {
			zap.S().Error("INDReadRedis ParseRedisLockVal Err. err:= %w", err)
			return "", err
		}

		if strStatus == lock.RedisWriteLocked && retryCount < 10 {
			retryCount++
			goto Retry
		}
		return "", err
	}

	cmd = client.Rdb.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			zap.S().Info("key does not exists. key:= %v", key)
			return "", nil
		}
		return "", cmd.Err()
	}

	return cmd.Val(), nil
}

// 直接的读取
func DReadRedis(client Client, key string) (string, error) {
	//本地超时
	var ctx, cancel = context.WithTimeout(context.Background(), time.Duration(lock.RedisTTL)*time.Second)
	defer cancel()

	cmd := client.Rdb.Get(ctx, key)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			zap.S().Info("key does not exists. key:= %v", key)
			return "", nil
		}
		return "", cmd.Err()
	}

	return cmd.Val(), nil
}
