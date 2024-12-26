package pgsql

const (
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultHost     = "localhost"
	defaultPort     = 5432
)

type PgsqlCredentials struct {
	User     string
	Password string
	Host     string
	Port     uint16
	Database string
}

func NewPgsqlCredentials(user string, password string, host string, port uint16, database string) PgsqlCredentials {
	return PgsqlCredentials{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

func NewDefaultPgsqlCredentials(database string) PgsqlCredentials {
	return PgsqlCredentials{
		User:     defaultUser,
		Password: defaultPassword,
		Host:     defaultHost,
		Port:     defaultPort,
		Database: database,
	}
}

type PgsqlSSLMode string

func (s PgsqlSSLMode) String() string {
	return string(s)
}

const (
	DisableMode    PgsqlSSLMode = "disable"
	AllowMode      PgsqlSSLMode = "allow"
	PreferMode     PgsqlSSLMode = "prefer"
	RequireMode    PgsqlSSLMode = "require"
	VerifyCAMode   PgsqlSSLMode = "verify-ca"
	VerifyFullMode PgsqlSSLMode = "verify-full"
)

type PgsqlClientOptionsFunc func(co *PgsqlClientOptions)

type PgsqlClientOptions struct {
	Credentials    PgsqlCredentials
	MaxConnections int
	ConnIdle       int
	MaxLifetime    int
	SSLMode        string
}

func NewDefaultClientOptions(credentials PgsqlCredentials) *PgsqlClientOptions {
	return &PgsqlClientOptions{
		Credentials:    credentials,
		MaxConnections: 20,
		ConnIdle:       50,
		MaxLifetime:    3,
		SSLMode:        VerifyFullMode.String(),
	}
}

func (pco *PgsqlClientOptions) apply(options ...PgsqlClientOptionsFunc) *PgsqlClientOptions {
	for _, opt := range options {
		opt(pco)
	}
	return pco
}

func WithMaxConnections(maxConnections int) PgsqlClientOptionsFunc {
	return func(co *PgsqlClientOptions) {
		co.MaxConnections = maxConnections
	}
}

func WithConnIdle(connIdle int) PgsqlClientOptionsFunc {
	return func(co *PgsqlClientOptions) {
		co.ConnIdle = connIdle
	}
}

func WithMaxLifetime(maxLifetime int) PgsqlClientOptionsFunc {
	return func(co *PgsqlClientOptions) {
		co.MaxLifetime = maxLifetime
	}
}

func WithSSLMode(mode PgsqlSSLMode) PgsqlClientOptionsFunc {
	return func(co *PgsqlClientOptions) {
		co.SSLMode = mode.String()
	}
}
