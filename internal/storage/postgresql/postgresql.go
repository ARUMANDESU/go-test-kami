package postgresql

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations
var MigrationsFs embed.FS

type Storage struct {
	Pool *pgxpool.Pool
}

func (s Storage) Close() {
	s.Pool.Close()
}

func NewStorage(ctx context.Context, dbURL string) (Storage, error) {

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to ping database: %w", err)
	}

	err = applyMigrations(dbURL, "reservations", "migrations", MigrationsFs)
	if err != nil {
		return Storage{}, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return Storage{Pool: dbpool}, nil
}

func applyMigrations(connStr, dbName, migrationTableName string, migrationsFs embed.FS) error {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: migrationTableName,
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
