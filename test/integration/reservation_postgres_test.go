package integration

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ARUMANDESU/go-test-kami/internal/storage/postgresql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start a PostgreSQL container
	var err error
	pgContainer, err = CreatePostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}

	// Apply migrations to the template database
	err = applyMigrations(pgContainer.ConnectionString, "template_db", "migration", postgresql.MigrationsFs)
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	// Apply mock data migrations
	err = applyMigrations(pgContainer.ConnectionString, "test_db", "mock_data", migrationsFS)
	if err != nil {
		log.Fatalf("failed to apply mock data migrations: %v", err)
	}

	code := m.Run()

	pgContainer.Terminate(ctx)

	os.Exit(code)
}

func TestReservationServise_GetRoomReservations_HappyPath(t *testing.T) {
	suite := NewReservationSuite(t)

	roomID := "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f"

	// When
	reservations, err := suite.service.GetRoomReservations(context.Background(), roomID)

	// Then
	require.NoError(t, err)

	assert.Len(t, reservations, 3)
}
