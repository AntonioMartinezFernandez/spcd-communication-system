package retry

import (
	"context"
	"log/slog"
	"time"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"

	"github.com/cenkalti/backoff/v3"
)

type RetryConfig struct {
	// MaxRetries is the maximum number of retries. If 0, retries are disabled.
	MaxRetries int
	// InitialInterval is the initial interval between retries. The interval will be increased exponentially.
	InitialInterval time.Duration
	// MaxInterval is the maximum interval between retries. Disabled if 0. Interval won't be greater than MaxInterval.
	MaxInterval time.Duration
	// Multiplier is the factor by which the interval will be increased on each retry.
	Multiplier float64
	// MaxElapsedTime is the maximum elapsed time between retries. Disabled if 0.
	MaxElapsedTime time.Duration
	// RandomizationFactor is the factor by which the backoff duration is randomized. Disabled if 0.
	RandomizationFactor float64

	// OnRetryScapeHook [optional] is called after each retry.
	// The first argument is the retry number, the second is the delay before the next retry and the third is the error.
	// If the returning value is true, the retry will be skipped.
	OnRetryScapeHook func(retryNum int, delay time.Duration, err error) bool

	Logger logger.Logger
}

func RetryBackoff(ctx context.Context, config RetryConfig, callback func() (interface{}, error)) (interface{}, error) {
	expoBackoff := backoff.NewExponentialBackOff()
	expoBackoff.InitialInterval = config.InitialInterval
	expoBackoff.MaxInterval = config.MaxInterval
	expoBackoff.Multiplier = config.Multiplier
	expoBackoff.MaxElapsedTime = config.MaxElapsedTime
	expoBackoff.RandomizationFactor = config.RandomizationFactor

	res, err := callback()
	if err == nil {
		return res, nil
	}

	if config.MaxElapsedTime > 0 {
		_, cancel := context.WithTimeout(ctx, config.MaxElapsedTime)
		defer cancel()
	}

	retryNum := 1
	expoBackoff.Reset()
retryLoop:
	for {
		timeToWait := expoBackoff.NextBackOff()
		select {
		case <-ctx.Done():
			return nil, err
		case <-time.After(timeToWait):
			// Let's continue...
		}

		if config.OnRetryScapeHook != nil && config.OnRetryScapeHook(retryNum, timeToWait, err) {
			break retryLoop
		}

		res, err := callback()
		if err == nil {
			return res, nil
		}

		if config.Logger != nil {
			config.Logger.Warn(
				ctx,
				"error retrying operation",
				slog.Int("retry_number", retryNum),
				slog.Int("max_retries", config.MaxRetries),
				slog.Duration("wait_time", timeToWait),
				slog.Duration("elapsed_time", expoBackoff.GetElapsedTime()),
				slog.String("error", err.Error()),
			)
		}

		retryNum++
		if retryNum > config.MaxRetries {
			break retryLoop
		}
	}

	return res, err
}
