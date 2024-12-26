package redis

/*

	REDIS-CLEANER

		- Redis record deletion tool, especially useful for deleting records with too high TTLs
		- Example of use:
				rcConfig := redis.NewRedisCleanerConfig("message:duplicator:guardian:*:*", 600, redis.WithBatchSize(100), redis.WithReport())
				redisCleaner := redis.NewRedisCleaner(redisClient, logger, rcConfig)
				redisCleaner.Run(ctx)

*/

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"

	"github.com/redis/go-redis/v9"
)

const RedisCleanerRunningKey = "redis-cleaner:is-running"

type RedisCleaner struct {
	redisClient *redis.Client
	logger      logger.Logger
	config      *RedisCleanerConfig

	keysToCheckChan  chan string
	keysToDeleteChan chan string

	mut       *sync.Mutex
	wg        *sync.WaitGroup
	startTime time.Time

	deleted int

	subCtxCancelFunc context.CancelFunc
}

func NewRedisCleaner(redisClient *redis.Client, logger logger.Logger, config *RedisCleanerConfig) *RedisCleaner {
	return &RedisCleaner{
		redisClient: redisClient,
		logger:      logger,
		config:      config,

		keysToCheckChan:  make(chan string),
		keysToDeleteChan: make(chan string),

		mut:       &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		startTime: time.Now(),

		deleted: 0,
	}
}

func (rc *RedisCleaner) Run(ctx context.Context) error {
	subCtx, cancelFunc := context.WithCancel(ctx)
	rc.mut.Lock()
	rc.subCtxCancelFunc = cancelFunc
	rc.mut.Unlock()

	newCtx := context.Background()

	defer func() {
		cancelFunc()
		_ = rc.setAsStopped(newCtx)
		if rc.config.showReport {
			rc.logger.Info(
				ctx,
				"redis-cleaner report",
				slog.Int("deleted_keys", rc.deleted),
				slog.Int64("duration_in_millis", time.Duration(time.Since(rc.startTime)).Abs().Milliseconds()),
			)
		}
		rc.logger.Info(ctx, "redis-cleaner stopped")
	}()

	// Check if the cleaner is already running
	isRunning, rErr := rc.isRunning(subCtx)
	if rErr != nil {
		rc.logger.Error(ctx, "error getting redis-cleaner running status", slog.String("error", rErr.Error()))
		return rErr
	}
	if isRunning {
		rc.logger.Info(ctx, "redis-cleaner is already running")
		return nil
	}

	srErr := rc.setAsRunning(subCtx)
	if srErr != nil {
		rc.logger.Error(ctx, "error setting redis-cleaner running status", slog.String("error", srErr.Error()))
	}

	// Start redis cleaner
	rc.logger.Info(ctx, "starting redis-cleaner...")

	rc.wg.Add(1)
	go rc.deleteKeys(newCtx)

	rc.wg.Add(1)
	go rc.checkKeys(newCtx)

	rc.wg.Add(1)
	go rc.findKeys(subCtx)

	// Wait for all goroutines to finish
	rc.wg.Wait()
	return nil
}

func (rc *RedisCleaner) Stop() {
	rc.mut.Lock()
	rc.subCtxCancelFunc()
	rc.mut.Unlock()
}

func (rc *RedisCleaner) findKeys(ctx context.Context) {
	defer func() {
		rc.wg.Done()
		close(rc.keysToCheckChan)
	}()

	scanIterator := rc.redisClient.Scan(ctx, 0, rc.config.keyRegex, 0).Iterator()
	for scanIterator.Next(ctx) {
		select {
		case <-ctx.Done():
			return
		default:
			rc.keysToCheckChan <- scanIterator.Val()
			if err := scanIterator.Err(); err != nil {
				rc.logger.Error(ctx, "redis-cleaner error iterating over keys", slog.String("error", err.Error()))
				rc.Stop()
			}
		}
	}
}

func (rc *RedisCleaner) checkKeys(ctx context.Context) {
	defer func() {
		close(rc.keysToDeleteChan)
		rc.wg.Done()
	}()

	for key := range rc.keysToCheckChan {
		keyTtlDuration, err := rc.redisClient.TTL(ctx, key).Result()
		if err != nil {
			rc.logger.Error(ctx, "error while redis-cleaner getting TTL", slog.String("error", err.Error()))
			rc.Stop()
		}
		if keyTtlDuration > time.Duration(rc.config.minTtlInSeconds)*time.Second {
			rc.keysToDeleteChan <- key
		}
	}
}

func (rc *RedisCleaner) deleteKeys(ctx context.Context) {
	defer func() {
		rc.wg.Done()
	}()

	keyBatch := make([]string, 0, rc.config.batchSize)

	for key := range rc.keysToDeleteChan {
		if len(keyBatch) >= rc.config.batchSize {
			rc.delete(ctx, keyBatch)
			keyBatch = keyBatch[:0]
		}
		keyBatch = append(keyBatch, key)
	}

	if len(keyBatch) > 0 {
		rc.delete(ctx, keyBatch)
		keyBatch = nil
	}
}

func (rc *RedisCleaner) delete(ctx context.Context, keyBatch []string) {
	err := rc.redisClient.Del(ctx, keyBatch...).Err()
	if err != nil {
		rc.logger.Error(ctx, "error while redis-cleaner deleting keys", slog.String("error", err.Error()))
		rc.Stop()
	}
	rc.mut.Lock()
	rc.deleted += len(keyBatch)
	rc.mut.Unlock()
}

func (rc *RedisCleaner) isRunning(ctx context.Context) (bool, error) {
	isRunning, err := rc.redisClient.Exists(ctx, RedisCleanerRunningKey).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}
	return isRunning == 1, nil
}

func (rc *RedisCleaner) setAsRunning(ctx context.Context) error {
	return rc.redisClient.Set(ctx, RedisCleanerRunningKey, "1", 0).Err()
}

func (rc *RedisCleaner) setAsStopped(ctx context.Context) error {
	return rc.redisClient.Del(ctx, RedisCleanerRunningKey).Err()
}
