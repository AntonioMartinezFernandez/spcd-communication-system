package di

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/AntonioMartinezFernandez/services/iot-devices/configs"

	amf_bus "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus"
	amf_command_bus "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/command"
	amf_query_bus "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/bus/query"
	amf_sync "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/distributed-sync"
	amf_json_schema "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/json-schema"
	amf_logger "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/logger"
	amf_observability "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/observability"
	amf_redis "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/redis"
	amf_sqldb "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/sqldb"
	amf_pgsql "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/sqldb/pgsql"
	amf_utils "github.com/AntonioMartinezFernandez/services/iot-devices/pkg/utils"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type CommonServices struct {
	Config      configs.Config
	Environment configs.Environment

	Logger                 amf_logger.Logger
	JsonSchemaValidator    *amf_json_schema.JsonSchemaValidator
	Observability          *amf_observability.OtelObservability
	RedisClient            *redis.Client
	DatabaseConnectionPool amf_sqldb.ConnectionPool
	DatabaseMigrator       amf_sqldb.Migrator
	DistributedMutex       amf_sync.MutexService
	UlidProvider           amf_utils.UlidProvider
	UuidProvider           amf_utils.UuidProvider
	TimeProvider           amf_utils.DateTimeProvider
	CommandBus             *amf_command_bus.CommandBus
	QueryBus               *amf_query_bus.QueryBus
}

func InitCommonServices(ctx context.Context) *CommonServices {
	config := initConfig()
	environment := configs.NewEnvironmentFromRawEnvVar(config.AppEnv)
	logger := amf_logger.NewOtelInstrumentalizedLogger(config.LogLevel)
	jsonSchemaValidator := amf_json_schema.NewJsonSchemaValidator(config.JsonSchemaBasePath)
	redisClient := amf_redis.NewRedisClient(config.RedisHost, config.RedisPort)
	redisMutexService := amf_sync.NewRedisMutexService(redisClient, logger)
	ulidProvider := amf_utils.NewRandomUlidProvider()
	uuidProvider := amf_utils.NewRandomUuidProvider()
	timeProvider := amf_utils.NewSystemTimeProvider()
	commandBus := amf_command_bus.InitCommandBus(logger, redisMutexService)
	queryBus := amf_query_bus.InitQueryBus(logger)
	databasePool := initPgsqlDatabasePool(ctx, config, environment)
	databaseMigrator := amf_sqldb.NewSQLDatabaseMigrator(
		databasePool.Writer(),
		config.MigrationsPath,
		"migrations",
		amf_sqldb.PgSQLPlatform,
	)

	grpcConnection, grpcErr := amf_observability.InitGrpcConnInsecure(config.OtelGrpcHost, config.OtelGrpcPort)
	if grpcErr != nil {
		logger.Error(ctx, "error establishing grpc connection with OTEL collector")
	}
	otelObservability, obsErr := amf_observability.InitOpenTelemetryObservability(
		ctx,
		grpcConnection,
		config.AppServiceName,
		config.AppVersion,
	)
	if obsErr != nil {
		logger.Error(ctx, "error initializing OpenTelemetry observability")
	}

	return &CommonServices{
		Config:      config,
		Environment: environment,

		Logger:                 logger,
		JsonSchemaValidator:    &jsonSchemaValidator,
		Observability:          otelObservability,
		RedisClient:            redisClient,
		DatabaseConnectionPool: databasePool,
		DatabaseMigrator:       databaseMigrator,
		DistributedMutex:       redisMutexService,
		UlidProvider:           ulidProvider,
		UuidProvider:           uuidProvider,
		TimeProvider:           timeProvider,
		CommandBus:             commandBus,
		QueryBus:               queryBus,
	}
}

/* HELPERS */

func InitCommonServicesWithEnvFiles(envFiles ...string) *CommonServices {
	ctx := context.Background()
	err := godotenv.Overload(envFiles...)
	if err != nil {
		panic(err)
	}

	return InitCommonServices(ctx)
}

func RootContext() (context.Context, context.CancelFunc) {
	rootCtx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, os.Kill,
	)
	return rootCtx, cancel
}

func initConfig() configs.Config {
	return configs.LoadEnvConfig()
}

func initPgsqlDatabasePool(_ context.Context, cfg configs.Config, env configs.Environment) *amf_pgsql.PgsqlConnectionPool {
	writerCredentials := amf_pgsql.NewPgsqlCredentials(
		cfg.PgsqlUser,
		cfg.PgsqlPassword,
		cfg.PgsqlHost,
		cfg.PgsqlPort,
		cfg.PgsqlDatabase,
	)

	readerCredentials := amf_pgsql.NewPgsqlCredentials(
		cfg.PgsqlUser,
		cfg.PgsqlPassword,
		cfg.PgsqlHostReader,
		cfg.PgsqlPort,
		cfg.PgsqlDatabase,
	)

	opts := make([]amf_pgsql.PgsqlClientOptionsFunc, 0)
	if env.IsDevelopment() {
		opts = append(opts, amf_pgsql.WithSSLMode(amf_pgsql.DisableMode))
	}

	writer, err := amf_pgsql.NewWriter(writerCredentials, opts...)
	if err != nil {
		panic(err)
	}

	reader, err := amf_pgsql.NewReader(readerCredentials, opts...)
	if err != nil {
		panic(err)
	}

	pool, err := amf_pgsql.NewPgsqlConnectionPool(writer, reader)

	if err != nil {
		panic(err)
	}

	return pool
}

func registerQueryOrPanic(
	queryBus amf_query_bus.Bus,
	query amf_bus.Dto,
	handler amf_query_bus.QueryHandler,
) {
	if err := queryBus.RegisterQuery(query, handler); err != nil {
		panic(err)
	}
}

func registerCommandOrPanic(
	commandBus amf_command_bus.Bus,
	cmd amf_bus.Dto,
	handler amf_command_bus.CommandHandler,
) {
	if err := commandBus.RegisterCommand(cmd, handler); err != nil {
		panic(err)
	}
}

func databaseMigrationFunc(ctx context.Context, common *CommonServices) func() (interface{}, error) {
	return func() (interface{}, error) {
		migrationsExecuted, err := common.DatabaseMigrator.Up()
		if err != nil {
			return nil, err
		}

		common.Logger.Warn(ctx, fmt.Sprintf("Applied %d migrations!", migrationsExecuted))
		return migrationsExecuted, nil
	}
}
