package sqldb

type Migrator interface {
	Up() (int, error)
	Down() (int, error)
}

type Platform string

func (p Platform) String() string {
	return string(p)
}

const (
	PgSQLPlatform Platform = "postgres"
)
