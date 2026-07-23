//go:build !nosql

package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type sqlDatabase struct {
	db     *sql.DB
	config *Config
}

func (s *sqlDatabase) defaultContext() (context.Context, context.CancelFunc) {
	if s.config != nil && s.config.Timeout > 0 {
		return context.WithTimeout(context.Background(), s.config.Timeout)
	}
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (s *sqlDatabase) close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *sqlDatabase) ping() error {
	if s.db == nil {
		return fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.db.PingContext(ctx)
}

func (s *sqlDatabase) exec(query string, args ...any) (Result, error) {
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.execContext(ctx, query, args...)
}

func (s *sqlDatabase) execContext(ctx context.Context, query string, args ...any) (Result, error) {
	if s.db == nil {
		return Result{}, fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return Result{}, err
	}
	affected, _ := res.RowsAffected()
	lastID, _ := res.LastInsertId()
	return Result{RowsAffected: affected, LastInsertID: lastID}, nil
}

func (s *sqlDatabase) query(query string, args ...any) ([]Row, error) {
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.queryContext(ctx, query, args...)
}

func (s *sqlDatabase) queryContext(ctx context.Context, query string, args ...any) ([]Row, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []Row
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		row := make(Row)
		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (s *sqlDatabase) queryRow(query string, args ...any) (Row, error) {
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.queryRowContext(ctx, query, args...)
}

func (s *sqlDatabase) queryRowContext(ctx context.Context, query string, args ...any) (Row, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return Row{}, nil
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	row := make(Row)
	for i, col := range columns {
		row[col] = values[i]
	}
	return row, nil
}

func (s *sqlDatabase) begin() (Transaction, error) {
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.beginTx(ctx)
}

func (s *sqlDatabase) beginTx(ctx context.Context) (Transaction, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx: tx}, nil
}

func (s *sqlDatabase) migrate(migrations []Migration) error {
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.migrateContext(ctx, migrations)
}

func (s *sqlDatabase) migrateContext(ctx context.Context, migrations []Migration) error {
	if s.db == nil {
		return fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}

	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS _migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			down_sql TEXT,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	for _, migration := range migrations {
		var count int
		err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM _migrations WHERE version = ?", migration.Version).Scan(&count)
		if err != nil {
			return fmt.Errorf("check migration %d: %w", migration.Version, err)
		}
		if count > 0 {
			continue
		}

		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %d: %w", migration.Version, err)
		}

		if _, err := tx.ExecContext(ctx, migration.Up); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %d: %w", migration.Version, err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO _migrations (version, name, down_sql) VALUES (?, ?, ?)", migration.Version, migration.Name, migration.Down); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", migration.Version, err)
		}
	}

	return nil
}

func (s *sqlDatabase) rollback(version int) error {
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.rollbackContext(ctx, version)
}

func (s *sqlDatabase) rollbackContext(ctx context.Context, version int) error {
	if s.db == nil {
		return fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}

	var migrations []Migration
	rows, err := s.db.QueryContext(ctx, "SELECT version, name, down_sql FROM _migrations WHERE version > ? ORDER BY version DESC", version)
	if err != nil {
		return fmt.Errorf("query migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var migration Migration
		if err := rows.Scan(&migration.Version, &migration.Name, &migration.Down); err != nil {
			return err
		}
		migrations = append(migrations, migration)
	}

	for _, migration := range migrations {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin rollback %d: %w", migration.Version, err)
		}

		if migration.Down != "" {
			if _, err := tx.ExecContext(ctx, migration.Down); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("execute down migration %d (%s): %w", migration.Version, migration.Name, err)
			}
		}

		if _, err := tx.ExecContext(ctx, "DELETE FROM _migrations WHERE version = ?", migration.Version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("remove migration record %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit rollback %d: %w", migration.Version, err)
		}
	}

	return nil
}

func (s *sqlDatabase) healthCheck() error {
	if s.db == nil {
		return fmt.Errorf("database not connected; call Connect() with a valid config before performing operations")
	}
	ctx, cancel := s.defaultContext()
	defer cancel()
	return s.db.PingContext(ctx)
}

type sqlTx struct {
	tx *sql.Tx
}

func (t *sqlTx) Exec(query string, args ...any) (Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return t.ExecContext(ctx, query, args...)
}

func (t *sqlTx) ExecContext(ctx context.Context, query string, args ...any) (Result, error) {
	res, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return Result{}, err
	}
	affected, _ := res.RowsAffected()
	lastID, _ := res.LastInsertId()
	return Result{RowsAffected: affected, LastInsertID: lastID}, nil
}

func (t *sqlTx) Query(query string, args ...any) ([]Row, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return t.QueryContext(ctx, query, args...)
}

func (t *sqlTx) QueryContext(ctx context.Context, query string, args ...any) ([]Row, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []Row
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		row := make(Row)
		for i, col := range columns {
			row[col] = values[i]
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (t *sqlTx) Commit() error {
	return t.tx.Commit()
}

func (t *sqlTx) Rollback() error {
	return t.tx.Rollback()
}
