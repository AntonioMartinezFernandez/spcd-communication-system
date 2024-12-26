package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

const driver = "postgres"

func NewReader(
	credentials PgsqlCredentials,
	opts ...PgsqlClientOptionsFunc,
) (*sql.DB, error) {
	options := NewDefaultClientOptions(credentials)
	options.apply(opts...)

	srvAddress, err := buildPgsqlConnectionStringFromOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := sql.Open(driver, srvAddress)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(); err != nil {
		return nil, err
	}

	client.SetConnMaxLifetime(time.Duration(options.MaxLifetime) * time.Minute)
	client.SetMaxOpenConns(options.MaxConnections)
	client.SetMaxIdleConns(options.ConnIdle)

	return client, nil
}

func NewWriter(
	credentials PgsqlCredentials,
	opts ...PgsqlClientOptionsFunc,
) (*sql.DB, error) {
	options := NewDefaultClientOptions(credentials)
	options.apply(opts...)

	srvAddress, err := buildPgsqlConnectionStringFromOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := sql.Open(driver, srvAddress)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(); err != nil {
		return nil, err
	}

	client.SetConnMaxLifetime(time.Duration(options.MaxLifetime) * time.Minute)
	client.SetMaxOpenConns(options.MaxConnections)
	client.SetMaxIdleConns(options.ConnIdle)

	return client, nil
}

func buildPgsqlConnectionStringFromOptions(options *PgsqlClientOptions) (string, error) {
	rawAddress := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		driver,
		options.Credentials.User,
		options.Credentials.Password,
		options.Credentials.Host,
		options.Credentials.Port,
		options.Credentials.Database,
		options.SSLMode,
	)

	srvAddress, err := pq.ParseURL(rawAddress)
	if err != nil {
		return "", err
	}

	return srvAddress, nil
}
