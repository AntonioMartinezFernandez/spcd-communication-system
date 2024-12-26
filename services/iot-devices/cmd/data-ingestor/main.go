package main

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/AntonioMartinezFernandez/services/iot-devices/cmd/di"
)

func main() {
	// Initialize Dependencies
	ctx, cancel := di.RootContext()
	errorsChannel := make(chan error)
	var wg = sync.WaitGroup{}
	di := di.InitDataIngestorDi(ctx)
	defer func() {
		cancel()
	}()

	di.CommonServices.Logger.Info(
		ctx,
		"starting HTTP server...",
		slog.String("service", di.CommonServices.Config.AppServiceName),
		slog.String("version", di.CommonServices.Config.AppVersion),
	)

	// Migrate up the database
	di.RunDatabaseMigrations(ctx)

	// Start Http Server
	go func() {
		errorsChannel <- di.HttpServices.Router.ListenAndServe(
			fmt.Sprintf("%s:%s", di.CommonServices.Config.HttpHost, di.CommonServices.Config.HttpPort),
		)
	}()

	// Shutdown servers on SIGINT, SIGTERM or error
	select {
	case err := <-errorsChannel:
		di.ErrorShutdown(ctx, cancel, err)
	case <-ctx.Done():
		di.GracefulShutdown(ctx)
	}

	wg.Wait()
}
