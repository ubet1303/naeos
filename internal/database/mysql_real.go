//go:build !nosql

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type RealMySQL struct {
	*sqlDatabase
}

func NewRealMySQL() *RealMySQL {
	return &RealMySQL{sqlDatabase: &sqlDatabase{}}
}

func (m *RealMySQL) Name() string {
	return "mysql"
}

func (m *RealMySQL) Connect(config *Config) error {
	m.config = config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	if config.Timeout > 0 {
		db.SetConnMaxLifetime(config.Timeout)
	}

	maxOpen := 25
	if config.MaxOpenConns > 0 {
		maxOpen = config.MaxOpenConns
	}
	db.SetMaxOpenConns(maxOpen)

	maxIdle := 5
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

	m.db = db
	return nil
}

func (m *RealMySQL) Close() error {
	return m.close()
}

func (m *RealMySQL) Ping() error {
	return m.ping()
}

func (m *RealMySQL) Exec(query string, args ...any) (Result, error) {
	return m.exec(query, args...)
}

func (m *RealMySQL) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return m.execContext(ctx, query, args...)
}

func (m *RealMySQL) Query(query string, args ...any) ([]Row, error) {
	return m.query(query, args...)
}

func (m *RealMySQL) QueryContext(ctx context.Context, query string, args ...any) ([]Row, error) {
	return m.queryContext(ctx, query, args...)
}

func (m *RealMySQL) QueryRow(query string, args ...any) (Row, error) {
	return m.queryRow(query, args...)
}

func (m *RealMySQL) QueryRowContext(ctx context.Context, query string, args ...any) (Row, error) {
	return m.queryRowContext(ctx, query, args...)
}

func (m *RealMySQL) Begin() (Transaction, error) {
	return m.begin()
}

func (m *RealMySQL) BeginTx(ctx context.Context) (Transaction, error) {
	return m.beginTx(ctx)
}

func (m *RealMySQL) Migrate(migrations []Migration) error {
	return m.migrate(migrations)
}

func (m *RealMySQL) MigrateContext(ctx context.Context, migrations []Migration) error {
	return m.migrateContext(ctx, migrations)
}

func (m *RealMySQL) Rollback(version int) error {
	return m.rollback(version)
}

func (m *RealMySQL) RollbackContext(ctx context.Context, version int) error {
	return m.rollbackContext(ctx, version)
}

func (m *RealMySQL) HealthCheck() error {
	return m.healthCheck()
}

// Type aliases for backward compatibility.
type RealMySQLTx = sqlTx
