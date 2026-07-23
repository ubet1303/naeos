//go:build !nosql

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type RealSQLite struct {
	*sqlDatabase
}

func NewRealSQLite() *RealSQLite {
	return &RealSQLite{sqlDatabase: &sqlDatabase{}}
}

func (s *RealSQLite) Name() string {
	return "sqlite"
}

func (s *RealSQLite) Connect(config *Config) error {
	s.config = config
	dsn := config.Database
	if dsn == "" {
		dsn = ":memory:"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	if config.Timeout > 0 {
		db.SetConnMaxLifetime(config.Timeout)
	}

	maxOpen := 1
	if config.MaxOpenConns > 0 {
		maxOpen = config.MaxOpenConns
	}
	db.SetMaxOpenConns(maxOpen)

	maxIdle := 1
	if config.MaxIdleConns > 0 {
		maxIdle = config.MaxIdleConns
	}
	db.SetMaxIdleConns(maxIdle)

	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}
	if config.ConnMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	if dsn != ":memory:" {
		if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
			db.Close()
			return fmt.Errorf("set WAL mode: %w", err)
		}
	}

	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.ExecContext(ctx, "PRAGMA busy_timeout=5000"); err != nil {
		db.Close()
		return fmt.Errorf("set busy timeout: %w", err)
	}

	s.db = db
	return nil
}

func (s *RealSQLite) Close() error {
	return s.close()
}

func (s *RealSQLite) Ping() error {
	return s.ping()
}

func (s *RealSQLite) Exec(query string, args ...any) (Result, error) {
	return s.exec(query, args...)
}

func (s *RealSQLite) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return s.execContext(ctx, query, args...)
}

func (s *RealSQLite) Query(query string, args ...any) ([]Row, error) {
	return s.query(query, args...)
}

func (s *RealSQLite) QueryContext(ctx context.Context, query string, args ...any) ([]Row, error) {
	return s.queryContext(ctx, query, args...)
}

func (s *RealSQLite) QueryRow(query string, args ...any) (Row, error) {
	return s.queryRow(query, args...)
}

func (s *RealSQLite) QueryRowContext(ctx context.Context, query string, args ...any) (Row, error) {
	return s.queryRowContext(ctx, query, args...)
}

func (s *RealSQLite) Begin() (Transaction, error) {
	return s.begin()
}

func (s *RealSQLite) BeginTx(ctx context.Context) (Transaction, error) {
	return s.beginTx(ctx)
}

func (s *RealSQLite) Migrate(migrations []Migration) error {
	return s.migrate(migrations)
}

func (s *RealSQLite) MigrateContext(ctx context.Context, migrations []Migration) error {
	return s.migrateContext(ctx, migrations)
}

func (s *RealSQLite) Rollback(version int) error {
	return s.rollback(version)
}

func (s *RealSQLite) RollbackContext(ctx context.Context, version int) error {
	return s.rollbackContext(ctx, version)
}

func (s *RealSQLite) HealthCheck() error {
	return s.healthCheck()
}

// Type alias for backward compatibility.
type RealSQLiteTx = sqlTx
