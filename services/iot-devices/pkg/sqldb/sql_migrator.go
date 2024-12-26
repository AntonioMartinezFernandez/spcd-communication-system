package sqldb

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

type SQLDatabaseMigrator struct {
	client          *sql.DB
	migrationSet    *migrate.MigrationSet
	migrationSource migrate.MigrationSource
	platform        Platform
}

func NewPgsqlDatabaseMigrator(client *sql.DB, migrationsPath, migrationsTableName string) *SQLDatabaseMigrator {
	return NewSQLDatabaseMigrator(client, migrationsPath, migrationsTableName, PgSQLPlatform)
}

func NewSQLDatabaseMigrator(
	client *sql.DB,
	migrationsPath,
	migrationsTableName string,
	platform Platform,
) *SQLDatabaseMigrator {
	migrationsSource := &migrate.FileMigrationSource{Dir: migrationsPath}
	migrationsSet := &migrate.MigrationSet{TableName: migrationsTableName}

	return &SQLDatabaseMigrator{
		client:          client,
		migrationSet:    migrationsSet,
		migrationSource: migrationsSource,
		platform:        platform,
	}
}

func (pdm *SQLDatabaseMigrator) Up() (int, error) {
	return pdm.migrationSet.Exec(pdm.client, pdm.platform.String(), pdm.migrationSource, migrate.Up)
}

func (pdm *SQLDatabaseMigrator) Down() (int, error) {
	return pdm.migrationSet.Exec(pdm.client, pdm.platform.String(), pdm.migrationSource, migrate.Down)
}
