package postgres

import (
	"compressor/internal/config"
	"compressor/pkg/logger"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

//DB ...
type DB struct {
	Session *sql.DB
	Logger  logger.Logger
}

//New ...
func New(logger logger.Logger, cfg config.Config) *DB {
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		logger.Fatalf("failed open conn to db %s", err)

		return nil
	}

	return &DB{
		Session: db,
		Logger:  logger,
	}
}

// Close ...
func (d *DB) Close() error {
	if err := d.Session.Close(); err != nil {
		return errors.Wrap(err, "can't close db")
	}

	return nil
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}
