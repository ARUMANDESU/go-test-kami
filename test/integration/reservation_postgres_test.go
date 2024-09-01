package integration

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
	"github.com/ARUMANDESU/go-test-kami/internal/service/reservation"
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

func TestReservationService_GetRoomReservations(t *testing.T) {
	tests := []struct {
		name     string
		roomID   string
		expected int
	}{
		{
			name:     "room has 3 reservations",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
			expected: 3,
		},
		{
			name:     "room has 1 reservation",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6b",
			expected: 1,
		},
		{
			name:     "room has 0 reservations",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5e",
			expected: 0,
		},
		{
			name:     "room does not exist",
			roomID:   "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5d",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewReservationSuite(t)

			reservations, err := s.service.GetRoomReservations(context.Background(), tt.roomID)
			require.NoError(t, err)

			assert.Len(t, reservations, tt.expected)
		})
	}
}

func TestReservationService_ReserveRoom_Concurrently(t *testing.T) {
	s := NewReservationSuite(t)

	var errCount int
	var wg sync.WaitGroup
	concurrentReservations := 20
	dto := domain.ReservationCreateDTO{
		RoomID:    "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b5f",
		StartTime: time.Date(2024, 8, 31, 12, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 8, 31, 14, 0, 0, 0, time.UTC),
	}

	wg.Add(concurrentReservations)
	for i := 0; i < concurrentReservations; i++ {
		go func() {
			defer wg.Done()
			reservedRoom, err := s.service.ReserveRoom(context.Background(), dto)

			if errors.Is(err, reservation.ErrReservationConflict) {
				errCount++
			} else {
				require.NoError(t, err)

				assert.NotEmpty(t, reservedRoom.ID)
				assert.Equal(t, dto.RoomID, reservedRoom.RoomID.String())
				assert.Equal(t, dto.StartTime, reservedRoom.StartTime)
				assert.Equal(t, dto.EndTime, reservedRoom.EndTime)
			}
		}()
	}

	wg.Wait()

	// Only one reservation should succeed
	assert.Equal(t, concurrentReservations-1, errCount)

	// Check the number of reservations
	// The room should have 4 reservations because the test mock data contains 3 reservations
	reservations, err := s.service.GetRoomReservations(context.Background(), dto.RoomID)
	require.NoError(t, err)

	assert.Len(t, reservations, 4)
}
