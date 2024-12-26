package redis

type RedisCleanerConfig struct {
	keyRegex        string
	batchSize       int
	minTtlInSeconds int
	showReport      bool
}

func (rco *RedisCleanerConfig) apply(opts ...RedisCleanerConfigFunc) *RedisCleanerConfig {
	for _, opt := range opts {
		opt(rco)
	}
	return rco
}

type RedisCleanerConfigFunc func(rco *RedisCleanerConfig)

func NewRedisCleanerConfig(
	keyRegex string,
	minTtlInSeconds int,
	opts ...RedisCleanerConfigFunc,
) *RedisCleanerConfig {
	options := &RedisCleanerConfig{
		keyRegex:        keyRegex,
		minTtlInSeconds: minTtlInSeconds,
		batchSize:       50,
		showReport:      false,
	}

	options.apply(opts...)
	return options
}

func WithBatchSize(batchSize int) RedisCleanerConfigFunc {
	return func(rco *RedisCleanerConfig) {
		if batchSize < 1 {
			batchSize = 50
		}
		rco.batchSize = batchSize
	}
}

func WithReport() RedisCleanerConfigFunc {
	return func(rco *RedisCleanerConfig) {
		rco.showReport = true
	}
}
