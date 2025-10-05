package migrate

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // регистрирует драйвер "pgx" для database/sql
)

// Встраиваем ВСЕ .sql из папки migrations.
//
//go:embed migrations/*.sql
var fs embed.FS

func newMigrator(dbURL string) (*migrate.Migrate, *sql.DB, error) {
	// database/sql на драйвере "pgx" — нужен для golang-migrate
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("postgres.WithInstance: %w", err)
	}

	src, err := iofs.New(fs, "migrations")
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("iofs.New: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("migrate.NewWithInstance: %w", err)
	}
	return m, db, nil
}

// Up — применяет все новые up-миграции (если нечего применять — не ошибка).
func Up(dbURL string) error {
	m, db, err := newMigrator(dbURL)
	if err != nil {
		return err
	}
	defer func() {
		if m != nil {
			// Close() возвращает 2 ошибки — просто игнорируем
			m.Close()
		}
		if db != nil {
			_ = db.Close()
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// DownAll — откатить все миграции до нуля (осторожно!).
func DownAll(dbURL string) error {
	m, db, err := newMigrator(dbURL)
	if err != nil {
		return err
	}
	defer func() {
		if m != nil {
			m.Close()
		}
		if db != nil {
			_ = db.Close()
		}
	}()

	return m.Down()
}

// Steps — сдвиг на N шагов (N>0 — up, N<0 — down).
func Steps(dbURL string, n int) error {
	m, db, err := newMigrator(dbURL)
	if err != nil {
		return err
	}
	defer func() {
		if m != nil {
			m.Close()
		}
		if db != nil {
			_ = db.Close()
		}
	}()

	return m.Steps(n)
}
