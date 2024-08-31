package reservation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
	"github.com/ARUMANDESU/go-test-kami/internal/service/reservation/mocks"
	"github.com/ARUMANDESU/go-test-kami/internal/storage"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thejerf/slogassert"
)

// suite is a test suite for the reservation service.
// It contains the common setup and teardown logic for the tests.
// It also contains the common fields used by the tests.
type suite struct {
	logHandler   *slogassert.Handler
	mockProvider *mocks.Provider
	mockReserver *mocks.Reserver
	service      Service
}

// NewSuite creates a new test suite.
// It creates a new assert logger handler(slogassert), mock provider, mock reserver, and service.
func NewSuite(t *testing.T) suite {
	t.Helper()
	handler := slogassert.New(t, slog.LevelWarn, nil)
	logger := slog.New(handler)

	mockProvider := mocks.NewProvider(t)
	mockReserver := mocks.NewReserver(t)

	service := NewService(logger, mockProvider, mockReserver)

	return suite{
		logHandler:   handler,
		mockProvider: mockProvider,
		mockReserver: mockReserver,
		service:      service,
	}
}

func TestGetRoomReservation_HappyPath(t *testing.T) {
	tests := []struct {
		name                      string
		roomID                    string
		onGetRoomReservations     []domain.Reservation
		onGetRoomReservationError error
	}{
		{
			name:                      "get room reservations",
			roomID:                    uuid.Must(uuid.NewV7()).String(),
			onGetRoomReservations:     MockReservations(t, time.Now()),
			onGetRoomReservationError: nil,
		},
		{
			name:                      "no reservations",
			roomID:                    uuid.Must(uuid.NewV7()).String(),
			onGetRoomReservations:     []domain.Reservation{},
			onGetRoomReservationError: storage.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			s.mockProvider.On("GetRoomReservations", mock.Anything, tt.roomID).Return(tt.onGetRoomReservations, tt.onGetRoomReservationError)

			reservations, err := s.service.GetRoomReservations(context.Background(), tt.roomID)
			require.NoError(t, err)

			require.Equal(t, tt.onGetRoomReservations, reservations)

			s.mockProvider.AssertExpectations(t)
		})
	}
}

func TestGetRoomReservation_FailPath(t *testing.T) {
	tests := []struct {
		name                      string
		roomID                    string
		onGetRoomReservations     []domain.Reservation
		onGetRoomReservationError error
		expectedError             error
	}{
		{
			name:                      "invalide room id: empty",
			roomID:                    "",
			onGetRoomReservations:     nil,
			onGetRoomReservationError: nil,
			expectedError:             ErrInvalidArgument,
		},
		{
			name:                      "invalid room id, not uuid",
			roomID:                    "invalid-uuid",
			onGetRoomReservations:     nil,
			onGetRoomReservationError: nil,
			expectedError:             ErrInvalidArgument,
		},
		{
			name:                      "failed to get room reservations",
			roomID:                    uuid.Must(uuid.NewV7()).String(),
			onGetRoomReservations:     nil,
			onGetRoomReservationError: fmt.Errorf("unexpected db error"),
			expectedError:             ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.expectedError, ErrInvalidArgument) {
				s.mockProvider.On("GetRoomReservations", mock.Anything, tt.roomID).Return(tt.onGetRoomReservations, tt.onGetRoomReservationError)
				defer s.mockProvider.AssertExpectations(t)
			}

			_, err := s.service.GetRoomReservations(context.Background(), tt.roomID)
			require.Error(t, err)

			assert.ErrorIs(t, err, tt.expectedError)

		})
	}

}

func TestReserveRoom_HappyPath(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		name                       string
		dto                        domain.ReservationCreateDTO
		onGetRoomReservations      []domain.Reservation
		onGetRoomReservationsError error
	}{
		{
			name: "reserve room: no reservations yet",
			dto: domain.ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow,
				EndTime:   timeNow.Add(time.Hour),
			},
			onGetRoomReservations:      nil,
			onGetRoomReservationsError: storage.ErrNotFound,
		},
		{
			name: "reserve room: with reservations",
			dto: domain.ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow,
				EndTime:   timeNow.Add(time.Hour),
			},
			onGetRoomReservations: MockReservations(t, timeNow),
		},
		{
			name: "reserve room: with reservations",
			dto: domain.ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow.Add(time.Hour * time.Duration(5)),
				EndTime:   timeNow.Add(time.Hour * time.Duration(6)),
			},
			onGetRoomReservations: MockReservations(t, timeNow),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			s.mockProvider.On("GetRoomReservations", mock.Anything, tt.dto.RoomID).Return(tt.onGetRoomReservations, tt.onGetRoomReservationsError)

			s.mockReserver.On("ReserveRoom", mock.Anything, tt.dto).Return(domain.Reservation{
				ID:        uuid.Must(uuid.NewV7()),
				RoomID:    uuid.FromStringOrNil(tt.dto.RoomID),
				StartTime: tt.dto.StartTime,
				EndTime:   tt.dto.EndTime,
			}, nil)

			reservation, err := s.service.ReserveRoom(context.Background(), tt.dto)
			require.NoError(t, err)

			assert.Equal(t, tt.dto.RoomID, reservation.RoomID.String())
			assert.Equal(t, tt.dto.StartTime, reservation.StartTime)
			assert.Equal(t, tt.dto.EndTime, reservation.EndTime)

			s.mockProvider.AssertExpectations(t)
			s.mockReserver.AssertExpectations(t)
		})
	}
}

func TestReserveRoom_FailPath(t *testing.T) {
	timeNow := time.Now()

	tests := []struct {
		name                       string
		dto                        domain.ReservationCreateDTO
		onGetRoomReservations      []domain.Reservation
		onGetRoomReservationsError error
		onReserveError             error
		expectedError              error
		onReserve                  bool
	}{
		{
			name: "invalid reservation create dto",
			dto: domain.ReservationCreateDTO{
				RoomID:    "",
				StartTime: timeNow,
				EndTime:   timeNow.Add(time.Hour),
			},
			onGetRoomReservations:      nil,
			onGetRoomReservationsError: nil,
			onReserveError:             nil,
			expectedError:              ErrInvalidArgument,
		},
		{
			name: "failed to get room reservations",
			dto: domain.ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow,
				EndTime:   timeNow.Add(time.Hour),
			},
			onGetRoomReservations:      nil,
			onGetRoomReservationsError: fmt.Errorf("unexpected db error"),
			onReserveError:             nil,
			expectedError:              ErrInternal,
		},
		{
			name: "reservation conflict",
			dto: domain.ReservationCreateDTO{
				RoomID: uuid.Must(uuid.NewV7()).String(),
				// this time period overlaps with the existing mock reservations
				StartTime: timeNow.Add(time.Hour),
				EndTime:   timeNow.Add(time.Hour * 2),
			},
			onGetRoomReservations:      MockReservations(t, timeNow),
			onGetRoomReservationsError: nil,
			onReserveError:             nil,
			expectedError:              ErrReservationConflict,
		},
		{
			name: "failed to reserve room, concurrency issue",
			dto: domain.ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow,
				EndTime:   timeNow.Add(time.Hour),
			},
			onGetRoomReservations:      MockReservations(t, timeNow),
			onGetRoomReservationsError: nil,
			onReserveError:             storage.ErrResevationConflict,
			expectedError:              ErrReservationConflict,
			onReserve:                  true,
		},
		{
			name: "failed to reserve room, internal error",
			dto: domain.ReservationCreateDTO{
				RoomID:    uuid.Must(uuid.NewV7()).String(),
				StartTime: timeNow,
				EndTime:   timeNow.Add(time.Hour),
			},
			onGetRoomReservations:      MockReservations(t, timeNow),
			onGetRoomReservationsError: nil,
			onReserveError:             fmt.Errorf("unexpected db error"),
			expectedError:              ErrInternal,
			onReserve:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSuite(t)

			if !errors.Is(tt.expectedError, ErrInvalidArgument) {
				s.mockProvider.On("GetRoomReservations", mock.Anything, tt.dto.RoomID).Return(tt.onGetRoomReservations, tt.onGetRoomReservationsError)
				defer s.mockProvider.AssertExpectations(t)

				if tt.onReserve {
					s.mockReserver.On("ReserveRoom", mock.Anything, tt.dto).Return(domain.Reservation{}, tt.onReserveError)
					defer s.mockReserver.AssertExpectations(t)
				}
			}

			_, err := s.service.ReserveRoom(context.Background(), tt.dto)
			require.Error(t, err)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}

}

// MockReservations is a helper function to create a list of reservations for testing purposes
// It creates 4 reservations with 1 hour difference
//
//	now                | now + 1 hour       | now + 2 hours      | now + 3 hours      | now + 4 hours      | now + 5 hours
//	<====================================================================================================================>
//	free slot (1 hour) | reserved (1 hour)  | reserved (1 hour)  | reserved (1 hour)  | reserved (1 hour)  | free slots after
//	free slot : 1 hour (now + 1 hour)
func MockReservations(t *testing.T, now time.Time) []domain.Reservation {
	t.Helper()

	var reservations []domain.Reservation
	for i := 1; i < 5; i++ {
		reservation := domain.Reservation{
			ID:        uuid.Must(uuid.NewV7()),
			RoomID:    uuid.Must(uuid.NewV7()),
			StartTime: now.Add(time.Hour * time.Duration(i)),
			EndTime:   now.Add(time.Hour * time.Duration(i+1)),
		}
		reservations = append(reservations, reservation)
	}

	return reservations
}
