package retry

import (
	"context"
	"testing"
	"time"

	amf_logger_mocks "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SpecialError struct {
}

func (s SpecialError) Error() string {
	return "special error"
}

func TestRetryBackoff(t *testing.T) {
	ctx := context.Background()

	t.Run("Success execution at once", func(t *testing.T) {
		logger := amf_logger_mocks.NewLogger(t)

		// Test configuration
		config := RetryConfig{
			MaxRetries:          3,
			InitialInterval:     100 * time.Millisecond,
			MaxInterval:         1 * time.Second,
			Multiplier:          2,
			MaxElapsedTime:      3 * time.Second,
			RandomizationFactor: 0.1,
			Logger:              logger,
		}

		sucessfullyCallable := func() (interface{}, error) {
			return "mock result", nil
		}

		result, err := RetryBackoff(ctx, config, sucessfullyCallable)

		assert.NoError(t, err)
		assert.Equal(t, result, "mock result")
	})

	t.Run("Failing especial error caught by hook", func(t *testing.T) {
		logger := amf_logger_mocks.NewLogger(t)

		// Test configuration
		config := RetryConfig{
			MaxRetries:          3,
			InitialInterval:     100 * time.Millisecond,
			MaxInterval:         1 * time.Second,
			Multiplier:          2,
			MaxElapsedTime:      3 * time.Second,
			RandomizationFactor: 0.1,
			Logger:              logger,
			OnRetryScapeHook: func(retryNum int, timeToWait time.Duration, err error) bool {
				if _, ok := err.(*SpecialError); ok {
					return true
				}
				return false
			},
		}

		failCallback := func() (interface{}, error) {
			return nil, &SpecialError{}
		}

		result, err := RetryBackoff(ctx, config, failCallback)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "special error")
		assert.Nil(t, result)
	})

	t.Run("Failing error executed N times (given by config)", func(t *testing.T) {
		logger := amf_logger_mocks.NewLogger(t)

		maxRetries := 3
		// Test configuration
		config := RetryConfig{
			MaxRetries:      maxRetries,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     1 * time.Second,
			MaxElapsedTime:  3 * time.Second,
			Logger:          logger,
		}

		failCallback := func() (interface{}, error) {
			return nil, &SpecialError{}
		}

		for i := 0; i < maxRetries; i++ {
			logger.On(
				"Warn",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return().Once()
		}

		result, err := RetryBackoff(ctx, config, failCallback)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "special error")
		assert.Nil(t, result)
	})
}
