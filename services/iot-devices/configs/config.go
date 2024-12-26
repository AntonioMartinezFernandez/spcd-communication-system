package configs

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	AppServiceName string `env:"APP_SERVICE_NAME"`
	AppEnv         string `env:"APP_ENV"`
	AppVersion     string `env:"APP_VERSION"`

	LogLevel           string `env:"LOG_LEVEL"`
	JsonSchemaBasePath string `env:"JSON_SCHEMA_BASE_PATH"`

	HttpHost         string `env:"HTTP_HOST"`
	HttpPort         string `env:"HTTP_PORT"`
	HttpReadTimeout  int    `env:"HTTP_READ_TIMEOUT"`
	HttpWriteTimeout int    `env:"HTTP_WRITE_TIMEOUT"`

	PgsqlHost       string `env:"PGSQL_HOST"`
	PgsqlHostReader string `env:"PGSQL_HOST_READER"`
	PgsqlUser       string `env:"PGSQL_USER"`
	PgsqlPassword   string `env:"PGSQL_PASS"`
	PgsqlPort       uint16 `env:"PGSQL_PORT"`
	PgsqlDatabase   string `env:"PGSQL_DATABASE"`
	MigrationsPath  string `env:"MIGRATIONS_PATH"`

	RedisHost string `env:"REDIS_HOST"`
	RedisPort int    `env:"REDIS_PORT"`

	OtelGrpcHost string `env:"OTEL_GRPC_HOST"`
	OtelGrpcPort string `env:"OTEL_GRPC_PORT"`

	DynamicParametersFilePath string `env:"DYNAMIC_PARAMETERS_FILE_PATH"`
	DynamicParametersApiKeys  string `env:"DYNAMIC_PARAMETERS_API_KEYS"`
}

func LoadEnvConfig() Config {
	var config Config

	_ = godotenv.Load(".env")
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}

	return config
}
