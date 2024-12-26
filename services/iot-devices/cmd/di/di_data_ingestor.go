package di

import (
	"context"
	"log/slog"
)

type DataIngestorDi struct {
	CommonServices           *CommonServices
	HttpServices             *HttpServices
	SystemServices           *SystemServices
	DynamicParameterServices *DynamicParameterServices
}

func InitDataIngestorDi(ctx context.Context) *DataIngestorDi {
	commonServices := InitCommonServices(ctx)
	httpServices := InitHttpServices(commonServices)
	systemServices := InitSystemServices(commonServices, httpServices)
	dynamicParameterServices := InitDynamicParameterServices(commonServices, httpServices)

	return &DataIngestorDi{
		CommonServices:           commonServices,
		HttpServices:             httpServices,
		SystemServices:           systemServices,
		DynamicParameterServices: dynamicParameterServices,
	}
}

func (iod *DataIngestorDi) RunDatabaseMigrations(ctx context.Context) {
	migrationFunc := databaseMigrationFunc(ctx, iod.CommonServices)
	_, err := iod.CommonServices.DistributedMutex.Mutex(ctx, "spcd_data_ingestor_migrations", migrationFunc)
	if err != nil {
		panic(err)
	}
}

func (iod *DataIngestorDi) ErrorShutdown(ctx context.Context, cancel context.CancelFunc, err error) {
	defer cancel()
	if err == nil {
		return
	}

	iod.HttpServices.Router.Shutdown(ctx)

	iod.CommonServices.Logger.Error(
		ctx,
		"error on starting server",
		slog.String("service", iod.CommonServices.Config.AppServiceName),
		slog.String("version", iod.CommonServices.Config.AppVersion),
		slog.String("error", err.Error()),
	)
}

func (iod *DataIngestorDi) GracefulShutdown(ctx context.Context) {
	iod.HttpServices.Router.Shutdown(ctx)

	iod.CommonServices.Logger.Info(
		ctx,
		"server stopped",
		slog.String("service", iod.CommonServices.Config.AppServiceName),
		slog.String("version", iod.CommonServices.Config.AppVersion),
	)
}
