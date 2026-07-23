//go:build !nosql

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type RealPostgreSQL struct {
	*sqlDatabase
}

func NewRealPostgreSQL() *RealPostgreSQL {
	return &RealPostgreSQL{sqlDatabase: &sqlDatabase{}}
}

func (p *RealPostgreSQL) Name() string {
	return "postgresql"
}

func (p *RealPostgreSQL) Connect(config *Config) error {
	p.config = config
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, config.SSLMode)

	db, err := sql.Open("pgx", dsn)
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

	p.db = db
	return nil
}

func (p *RealPostgreSQL) Close() error {
	return p.close()
}

func (p *RealPostgreSQL) Ping() error {
	return p.ping()
}

func (p *RealPostgreSQL) Exec(query string, args ...any) (Result, error) {
	return p.exec(query, args...)
}

func (p *RealPostgreSQL) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	return p.execContext(ctx, query, args...)
}

func (p *RealPostgreSQL) Query(query string, args ...any) ([]Row, error) {
	return p.query(query, args...)
}

func (p *RealPostgreSQL) QueryContext(ctx context.Context, query string, args ...any) ([]Row, error) {
	return p.queryContext(ctx, query, args...)
}

func (p *RealPostgreSQL) QueryRow(query string, args ...any) (Row, error) {
	return p.queryRow(query, args...)
}

func (p *RealPostgreSQL) QueryRowContext(ctx context.Context, query string, args ...any) (Row, error) {
	return p.queryRowContext(ctx, query, args...)
}

func (p *RealPostgreSQL) Begin() (Transaction, error) {
	return p.begin()
}

func (p *RealPostgreSQL) BeginTx(ctx context.Context) (Transaction, error) {
	return p.beginTx(ctx)
}

func (p *RealPostgreSQL) Migrate(migrations []Migration) error {
	return p.migrate(migrations)
}

func (p *RealPostgreSQL) MigrateContext(ctx context.Context, migrations []Migration) error {
	return p.migrateContext(ctx, migrations)
}

func (p *RealPostgreSQL) Rollback(version int) error {
	return p.rollback(version)
}

func (p *RealPostgreSQL) RollbackContext(ctx context.Context, version int) error {
	return p.rollbackContext(ctx, version)
}

func (p *RealPostgreSQL) HealthCheck() error {
	return p.healthCheck()
}

// Type alias for backward compatibility.
type RealPostgreSQLTx = sqlTx
