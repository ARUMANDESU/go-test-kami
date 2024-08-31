package reservation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
	"github.com/ARUMANDESU/go-test-kami/internal/storage"
	"github.com/ARUMANDESU/go-test-kami/pkg/logger"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofrs/uuid"
)

//go:generate mockery --name Provider --output mocks --case underscore
type Provider interface {
	// GetRoomReservations returns all reservations for a given room.
	// If no reservations are found, the function returns an `ErrNotFound` error.
	GetRoomReservations(ctx context.Context, roomID string) ([]domain.Reservation, error)
}

//go:generate mockery --name Reserver --output mocks --case underscore
type Reserver interface {
	// ReserveRoom reserves a room for a given time range.
	// It returns the reservation if the room is available, otherwise it returns an error.
	// If the reservation conflicts with another reservation, the function returns an `ErrResevationConflict` error.
	ReserveRoom(ctx context.Context, reservation domain.Reservation) (domain.Reservation, error)
}

type Service struct {
	log      *slog.Logger
	provider Provider
	reserver Reserver
}

func NewService(log *slog.Logger, provider Provider, reserver Reserver) Service {
	return Service{
		log:      log,
		provider: provider,
		reserver: reserver,
	}
}

func (s Service) ReserveRoom(ctx context.Context, dto domain.ReservationCreateDTO) (domain.Reservation, error) {
	const op = "reservation.Service.ReserveRoom"
	log := s.log.With("op", op)

	err := domain.ValidateReservationCreateDTO(dto)
	if err != nil {
		log.Debug("invalid reservation create dto", logger.Err(err))
		return domain.Reservation{}, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	rid, err := uuid.NewV7()
	if err != nil {
		log.Error("failed to generate reservation id", logger.Err(err))
		return domain.Reservation{}, ErrInternal
	}

	newReservation := domain.Reservation{
		ID:        rid,
		RoomID:    uuid.FromStringOrNil(dto.RoomID),
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
	}

	// check if reservation time is available
	reservations, err := s.provider.GetRoomReservations(ctx, dto.RoomID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			// if no reservations are found, continue to reserve the room
		default:
			log.Error("failed to get room reservations", logger.Err(err))
			return domain.Reservation{}, ErrInternal
		}
	}

	for _, r := range reservations {
		if newReservation.OverlapsWith(r) {
			return domain.Reservation{}, ErrReservationConflict
		}
	}

	reservation, err := s.reserver.ReserveRoom(ctx, newReservation)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrResevationConflict):
			return domain.Reservation{}, ErrReservationConflict
		default:
			log.Error("failed to reserve room", logger.Err(err))
			return domain.Reservation{}, ErrInternal
		}

	}

	return reservation, nil
}

func (s Service) GetRoomReservations(ctx context.Context, roomID string) ([]domain.Reservation, error) {
	const op = "reservation.Service.GetRoomReservations"
	log := s.log.With("op", op)

	err := validation.Validate(roomID, validation.Required, is.UUID)
	if err != nil {
		log.Debug("invalid room id", logger.Err(err))
		return nil, fmt.Errorf("%w: %w", ErrInvalidArgument, err)
	}

	reservations, err := s.provider.GetRoomReservations(ctx, roomID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []domain.Reservation{}, nil
		}
		log.Error("failed to get room reservations", logger.Err(err))
		return nil, ErrInternal
	}

	return reservations, nil
}
