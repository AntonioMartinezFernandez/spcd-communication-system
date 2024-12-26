package pgsql

import (
	"database/sql"

	"github.com/AntonioMartinezFernandez/services/iot-devices/pkg/sqldb"
)

type PgsqlConnectionPool struct {
	writer *sql.DB
	reader *sql.DB
}

func NewPgsqlConnectionPoolFromCredentials(
	writerCredentials,
	readerCredentials PgsqlCredentials,
) (*PgsqlConnectionPool, error) {
	writer, err := NewWriter(writerCredentials)
	if err != nil {
		return nil, err
	}

	reader, err := NewReader(readerCredentials)
	if err != nil {
		return nil, err
	}

	return NewPgsqlConnectionPool(writer, reader)
}

func NewPgsqlConnectionPool(writer, reader *sql.DB) (*PgsqlConnectionPool, error) {
	if writer == nil || reader == nil {
		return nil, sqldb.NewInvalidPoolConfigProvided(driver)
	}

	return &PgsqlConnectionPool{
		writer: writer,
		reader: reader,
	}, nil
}

func WithWriterOnly(writer *sql.DB) (*PgsqlConnectionPool, error) {
	if writer == nil {
		return nil, sqldb.NewInvalidPoolConfigProvided(driver)
	}

	return &PgsqlConnectionPool{
		writer: writer,
		reader: nil,
	}, nil
}

func (p *PgsqlConnectionPool) Writer() *sql.DB {
	return p.writer
}

func (p *PgsqlConnectionPool) Reader() *sql.DB {
	if nil == p.reader {
		return p.writer
	}

	return p.reader
}
