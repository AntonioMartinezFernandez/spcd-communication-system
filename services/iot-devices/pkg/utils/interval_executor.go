package utils

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
)

type ExecutorFunc func(ctx context.Context) error

func IntervalExecutor(
	ctx context.Context,
	handler ExecutorFunc,
	logger logger.Logger,
	ticker *time.Ticker,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	err := handler(ctx)
	if err != nil {
		logger.Error(ctx, "[interval executor] Error executing command", slog.String("error", err.Error()))
	}

	for {
		select {
		case <-ctx.Done():
			logger.Warn(ctx, "[interval executor] Terminated Signal received. Exiting command")
			return
		case <-ticker.C:
			err := handler(ctx)

			if err != nil {
				logger.Error(ctx, "[interval executor] Error executing command", slog.String("error", err.Error()))
			}
		}
	}
}
