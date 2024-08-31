package integration

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"testing"

	"github.com/ARUMANDESU/go-test-kami/internal/service/reservation"
	"github.com/ARUMANDESU/go-test-kami/internal/storage/postgresql"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/thejerf/slogassert"
)

//go:embed migrations
var migrationsFS embed.FS

var pgContainer *PostgresContainer

type ReservationSuite struct {
	loggerHandler *slogassert.Handler
	storage       postgresql.Storage
	service       reservation.Service
}

func NewReservationSuite(t *testing.T) *ReservationSuite {
	t.Helper()

	loggerHandler := slogassert.New(t, slog.LevelDebug, nil)
	logger := slog.New(loggerHandler)
	// Create a new test database cloned from the template database
	storage, cleanup := setupTestDB(t)

	t.Cleanup(func() {
		cleanup()
	})

	return &ReservationSuite{
		loggerHandler: loggerHandler,
		storage:       storage,
		service: reservation.NewService(
			logger,
			storage,
			storage,
		),
	}
}

func setupTestDB(t *testing.T) (postgresql.Storage, func()) {
	ctx := context.Background()

	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	dbName := "testdb_" + uuid.New().String()
	dsn := fmt.Sprintf("postgres://postgres:postgres@%s:%s/%s?sslmode=disable", host, port.Port(), dbName)

	// Connect to the default database to create a new test database
	defaultDB, err := pgxpool.New(ctx, fmt.Sprintf("postgres://postgres:postgres@%s:%s/template_db?sslmode=disable", host, port.Port()))
	if err != nil {
		t.Fatalf("failed to connect to default database: %v", err)
	}

	// Terminate all connections to the template_db
	_, err = defaultDB.Exec(ctx, "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'template_db' AND pid <> pg_backend_pid();")
	if err != nil {
		t.Fatalf("failed to terminate connections to template_db: %v", err)
	}

	_, err = defaultDB.Exec(ctx, fmt.Sprintf("CREATE DATABASE \"%s\" TEMPLATE template_db", dbName))
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	testDB, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return postgresql.Storage{Pool: testDB}, func() {
		testDB.Close()
		_, err := defaultDB.Exec(ctx, fmt.Sprintf("DROP DATABASE \"%s\"", dbName))
		if err != nil {
			t.Fatalf("failed to drop test database: %v", err)
		}

		defaultDB.Close()
	}
}

type PostgresContainer struct {
	testcontainers.Container
	ConnectionString string
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "template_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("postgres://postgres:postgres@%s:%s/template_db?sslmode=disable", host, port.Port())
	return &PostgresContainer{pgContainer, connStr}, nil
}

func applyMigrations(connStr, dbName, migrationTableName string, migrationsFs embed.FS) error {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{
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
