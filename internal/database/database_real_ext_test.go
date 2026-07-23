//go:build !nosql

package database

import (
	"context"
	"testing"
	"time"
)

func TestRealMySQLCloseNotConnected(t *testing.T) {
	db := NewRealMySQL()
	if err := db.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealPostgreSQLCloseNotConnected(t *testing.T) {
	db := NewRealPostgreSQL()
	if err := db.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealMySQLDefaultContextWithTimeout(t *testing.T) {
	db := NewRealMySQL()
	db.config = &Config{Timeout: 5 * time.Second}
	ctx, cancel := db.defaultContext()
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestRealMySQLDefaultContextWithoutTimeout(t *testing.T) {
	db := NewRealMySQL()
	db.config = &Config{}
	ctx, cancel := db.defaultContext()
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestRealMySQLDefaultContextNilConfig(t *testing.T) {
	db := NewRealMySQL()
	ctx, cancel := db.defaultContext()
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestRealPostgreSQLDefaultContextWithTimeout(t *testing.T) {
	db := NewRealPostgreSQL()
	db.config = &Config{Timeout: 5 * time.Second}
	ctx, cancel := db.defaultContext()
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestRealPostgreSQLDefaultContextWithoutTimeout(t *testing.T) {
	db := NewRealPostgreSQL()
	db.config = &Config{}
	ctx, cancel := db.defaultContext()
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestRealPostgreSQLDefaultContextNilConfig(t *testing.T) {
	db := NewRealPostgreSQL()
	ctx, cancel := db.defaultContext()
	defer cancel()
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestRealMySQLConnectErrorWithConnMaxLifetime(t *testing.T) {
	db := NewRealMySQL()
	err := db.Connect(&Config{
		Host:            "192.0.2.1",
		Port:            3306,
		User:            "root",
		Password:        "pass",
		Database:        "test",
		Timeout:         1 * time.Second,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealPostgreSQLConnectErrorWithConnMaxLifetime(t *testing.T) {
	db := NewRealPostgreSQL()
	err := db.Connect(&Config{
		Host:            "192.0.2.1",
		Port:            5432,
		User:            "postgres",
		Password:        "pass",
		Database:        "test",
		SSLMode:         "disable",
		Timeout:         1 * time.Second,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealMySQLCloseAfterConnectError(t *testing.T) {
	db := NewRealMySQL()
	_ = db.Connect(&Config{
		Host:     "192.0.2.1",
		Port:     3306,
		User:     "root",
		Database: "test",
		Timeout:  1 * time.Second,
	})
	if err := db.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealPostgreSQLCloseAfterConnectError(t *testing.T) {
	db := NewRealPostgreSQL()
	_ = db.Connect(&Config{
		Host:     "192.0.2.1",
		Port:     5432,
		User:     "postgres",
		Database: "test",
		SSLMode:  "disable",
		Timeout:  1 * time.Second,
	})
	if err := db.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRealSQLiteConnectWithFile(t *testing.T) {
	db := NewRealSQLite()
	tmpDir := t.TempDir()
	err := db.Connect(&Config{
		Database:        tmpDir + "/test.db",
		MaxOpenConns:    2,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		Timeout:         5 * time.Second,
	})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
}

func TestRealSQLiteConnectEmptyDSN(t *testing.T) {
	db := NewRealSQLite()
	err := db.Connect(&Config{Database: ""})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()
}

func TestRealSQLiteConnectBadPath(t *testing.T) {
	db := NewRealSQLite()
	err := db.Connect(&Config{Database: "/nonexistent/dir/test.db"})
	if err == nil {
		t.Error("expected error for bad path")
	}
}

func TestRealSQLiteQueryContextNoRows(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.Exec("CREATE TABLE empty_t (id INTEGER)")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	rows, err := db.QueryContext(context.Background(), "SELECT id FROM empty_t WHERE id = 999")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rows))
	}
}

func TestRealSQLiteQueryRowContextNoRows(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.Exec("CREATE TABLE empty_t (id INTEGER)")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	row, err := db.QueryRowContext(context.Background(), "SELECT id FROM empty_t WHERE id = 999")
	if err != nil {
		t.Fatalf("queryrow: %v", err)
	}
	if row == nil {
		t.Fatal("expected non-nil row")
	}
}

func TestRealSQLiteMigrateContextWithSkip(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	migrations := []Migration{
		{Version: 1, Name: "init", Up: "CREATE TABLE IF NOT EXISTS t1(id INT)", Down: "DROP TABLE IF EXISTS t1"},
		{Version: 2, Name: "add", Up: "CREATE TABLE IF NOT EXISTS t2(id INT)", Down: "DROP TABLE IF EXISTS t2"},
	}

	if err := db.MigrateContext(context.Background(), migrations); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if err := db.MigrateContext(context.Background(), migrations); err != nil {
		t.Fatalf("migrate again (should skip): %v", err)
	}
}

func TestRealSQLiteRollbackContextWithMigrations(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	migrations := []Migration{
		{Version: 1, Name: "v1", Up: "CREATE TABLE IF NOT EXISTS t1(id INT)", Down: "DROP TABLE IF EXISTS t1"},
		{Version: 2, Name: "v2", Up: "CREATE TABLE IF NOT EXISTS t2(id INT)", Down: "DROP TABLE IF EXISTS t2"},
	}

	if err := db.MigrateContext(context.Background(), migrations); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if err := db.RollbackContext(context.Background(), 0); err != nil {
		t.Fatalf("rollback: %v", err)
	}
}

func TestRealSQLiteTxExecAndQueryContext(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.Exec("CREATE TABLE tx_test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}

	_, err = tx.Exec("INSERT INTO tx_test (name) VALUES (?)", "tx-item")
	if err != nil {
		t.Fatalf("tx exec: %v", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO tx_test (name) VALUES (?)", "ctx-item")
	if err != nil {
		t.Fatalf("tx exec context: %v", err)
	}

	rows, err := tx.Query("SELECT name FROM tx_test")
	if err != nil {
		t.Fatalf("tx query: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}

	rows, err = tx.QueryContext(ctx, "SELECT name FROM tx_test")
	if err != nil {
		t.Fatalf("tx query context: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows from QueryContext, got %d", len(rows))
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}
}

func TestRealSQLiteTxRollback(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.Exec("CREATE TABLE tx_rb (id INTEGER)")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	tx, err := db.BeginTx(context.Background())
	if err != nil {
		t.Fatalf("begin: %v", err)
	}

	_, err = tx.Exec("INSERT INTO tx_rb VALUES (1)")
	if err != nil {
		t.Fatalf("tx exec: %v", err)
	}

	if err := tx.Rollback(); err != nil {
		t.Fatalf("rollback: %v", err)
	}

	rows, err := db.Query("SELECT id FROM tx_rb")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected 0 rows after rollback, got %d", len(rows))
	}
}

func TestRealSQLiteQueryContextError(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.QueryContext(context.Background(), "NOT VALID SQL")
	if err == nil {
		t.Error("expected error for bad SQL")
	}
}

func TestRealSQLiteQueryRowContextError(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.QueryRowContext(context.Background(), "NOT VALID SQL")
	if err == nil {
		t.Error("expected error for bad SQL")
	}
}

func TestRealSQLiteExecContextError(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.ExecContext(context.Background(), "NOT VALID SQL")
	if err == nil {
		t.Error("expected error for bad SQL")
	}
}

func TestRealSQLiteMigrateContextError(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	migrations := []Migration{
		{Version: 1, Name: "bad", Up: "NOT VALID SQL", Down: "SELECT 1"},
	}
	err := db.MigrateContext(context.Background(), migrations)
	if err == nil {
		t.Error("expected error for bad migration")
	}
}

func TestRealSQLiteRollbackContextError(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.Exec("CREATE TABLE _migrations (version INTEGER PRIMARY KEY, name TEXT, down_sql TEXT, applied_at TIMESTAMP)")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = db.Exec("INSERT INTO _migrations (version, name, down_sql) VALUES (1, 'v1', 'DROP TABLE bad')")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}

	err = db.RollbackContext(context.Background(), 0)
	if err == nil {
		t.Error("expected error for bad down migration")
	}
}

func TestRealSQLiteQueryContextNotConnected(t *testing.T) {
	db := NewRealSQLite()
	_, err := db.QueryContext(context.Background(), "SELECT 1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealSQLiteQueryRowContextNotConnected(t *testing.T) {
	db := NewRealSQLite()
	_, err := db.QueryRowContext(context.Background(), "SELECT 1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealSQLiteExecContextNotConnected(t *testing.T) {
	db := NewRealSQLite()
	_, err := db.ExecContext(context.Background(), "SELECT 1")
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealSQLiteMigrateContextNotConnected(t *testing.T) {
	db := NewRealSQLite()
	err := db.MigrateContext(context.Background(), nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealSQLiteRollbackContextNotConnected(t *testing.T) {
	db := NewRealSQLite()
	err := db.RollbackContext(context.Background(), 0)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRealSQLiteName(t *testing.T) {
	db := NewRealSQLite()
	if db.Name() != "sqlite" {
		t.Errorf("expected name 'sqlite', got %s", db.Name())
	}
}

func TestRealSQLiteBegin(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	_, err := db.Begin()
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
}

func TestRealSQLiteHealthCheck(t *testing.T) {
	db := NewRealSQLite()
	if err := db.Connect(&Config{Database: ":memory:"}); err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()

	if err := db.HealthCheck(); err != nil {
		t.Fatalf("health check: %v", err)
	}
}

func TestRealSQLiteHealthCheckNotConnected(t *testing.T) {
	db := NewRealSQLite()
	if err := db.HealthCheck(); err == nil {
		t.Error("expected error when not connected")
	}
}
