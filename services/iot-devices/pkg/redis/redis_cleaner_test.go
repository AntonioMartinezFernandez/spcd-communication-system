package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
	amf_redis "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/redis"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

// RedisCleanerTestSuite is a suite for testing RedisCleaner
type RedisCleanerTestSuite struct {
	suite.Suite
	redisClient *redis.Client
	logger      logger.Logger
	miniRedis   *miniredis.Miniredis
	ctx         context.Context
}

// SetupSuite is a setup function run before the suite starts
func (suite *RedisCleanerTestSuite) SetupSuite() {
	suite.logger = &logger.NullLogger{}

	// Create a miniRedis instance for testing
	miniRedis, err := miniredis.Run()
	suite.Require().NoError(err)
	suite.miniRedis = miniRedis

	// Create a Redis client with the miniRedis URL
	client := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	suite.redisClient = client
}

// SetupTest is a setup function run before each test in this suite
func (suite *RedisCleanerTestSuite) SetupTest() {
	suite.ctx = context.Background()

	// Set a redis key that SHOULD NOT be deleted
	err := suite.miniRedis.Set("test:key:not-expire", "not-expire")
	suite.Require().NoError(err)
	suite.miniRedis.SetTTL("test:key:not-expire", 1*time.Minute)

	// Create several redis keys that SHOULD be deleted
	for i := 1; i <= 1000; i++ {
		randomString := string(rune(i + 1))
		err = suite.miniRedis.Set("test:key:"+randomString, randomString)
		suite.Require().NoError(err)
		suite.miniRedis.SetTTL("test:key:"+randomString, 1*time.Hour)
	}
}

// TearDownSuite is a teardown function run after the suite finishes
func (suite *RedisCleanerTestSuite) TearDownSuite() {
	suite.miniRedis.Close()
	suite.redisClient.Close()
}

// TestRedisCleaner
func (suite *RedisCleanerTestSuite) TestRedisCleaner() {
	redisCleanerConfig := amf_redis.NewRedisCleanerConfig("test:key:*", 100)
	redisCleaner := amf_redis.NewRedisCleaner(suite.redisClient, suite.logger, redisCleanerConfig)
	err := redisCleaner.Run(suite.ctx)

	suite.NoError(err)
	suite.Equal(int64(0), suite.redisClient.Exists(suite.ctx, "redis-cleaner:is-running").Val())
	suite.Equal(int64(1), suite.redisClient.Exists(suite.ctx, "test:key:not-expire").Val())
	suite.Equal("not-expire", suite.redisClient.Get(suite.ctx, "test:key:not-expire").Val())
	suite.Equal(int64(0), suite.redisClient.Exists(suite.ctx, "test:key:1").Val())

	keys := suite.redisClient.Keys(suite.ctx, "test:key:*").Val()
	suite.Equal(int(1), len(keys))
}

func (suite *RedisCleanerTestSuite) TestRedisCleanerWhenContextIsCancelled() {
	redisCleanerConfig := amf_redis.NewRedisCleanerConfig("test:key:*", 100)
	redisCleaner := amf_redis.NewRedisCleaner(suite.redisClient, suite.logger, redisCleanerConfig)

	ctx, cancel := context.WithCancel(suite.ctx)
	go func() {
		<-time.After(15 * time.Millisecond)
		cancel()
	}()

	err := redisCleaner.Run(ctx)

	suite.NoError(err)
	suite.Equal(int64(0), suite.redisClient.Exists(suite.ctx, "redis-cleaner:is-running").Val())
	suite.Equal(int64(1), suite.redisClient.Exists(suite.ctx, "test:key:not-expire").Val())
	suite.Equal("not-expire", suite.redisClient.Get(suite.ctx, "test:key:not-expire").Val())
	suite.Equal(int64(0), suite.redisClient.Exists(suite.ctx, "test:key:1").Val())

	remainingKeys := suite.redisClient.Keys(suite.ctx, "test:key:*").Val()
	suite.Greater(len(remainingKeys), int(10))
	suite.Less(len(remainingKeys), int(1000))
}

func (suite *RedisCleanerTestSuite) TestRedisCleanerWhenIsStopped() {
	redisCleanerConfig := amf_redis.NewRedisCleanerConfig("test:key:*", 100)
	redisCleaner := amf_redis.NewRedisCleaner(suite.redisClient, suite.logger, redisCleanerConfig)

	go func() {
		<-time.After(15 * time.Millisecond)
		redisCleaner.Stop()
	}()

	err := redisCleaner.Run(suite.ctx)

	suite.NoError(err)
	suite.Equal(int64(0), suite.redisClient.Exists(suite.ctx, "redis-cleaner:is-running").Val())
	suite.Equal(int64(1), suite.redisClient.Exists(suite.ctx, "test:key:not-expire").Val())
	suite.Equal("not-expire", suite.redisClient.Get(suite.ctx, "test:key:not-expire").Val())
	suite.Equal(int64(0), suite.redisClient.Exists(suite.ctx, "test:key:1").Val())

	remainingKeys := suite.redisClient.Keys(suite.ctx, "test:key:*").Val()
	suite.Greater(len(remainingKeys), int(10))
	suite.Less(len(remainingKeys), int(1000))
}

func TestRedisCleanerSuite(t *testing.T) {
	suite.Run(t, new(RedisCleanerTestSuite))
}
