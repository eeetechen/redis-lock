package distributed_lock

const (
	RedisReadLocked uint8 = iota + 1
	RedisWriteLocked
)

const RedisTTL int64 = 15

const RetryMaxTime int = 10
