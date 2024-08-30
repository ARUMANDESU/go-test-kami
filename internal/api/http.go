package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
)

type ReservationService interface {
	ReserveRoom(ctx context.Context, dto domain.ReservationCreateDTO) (domain.Reservation, error)
	GetRoomReservations(ctx context.Context, roomID string) ([]domain.Reservation, error)
}

type API struct {
	log                *slog.Logger
	ReservationService ReservationService
}

func NewAPI(logger *slog.Logger, reservationService ReservationService) API {
	return API{
		log:                logger,
		ReservationService: reservationService,
	}
}

func (a API) HTTPServer(addr string) http.Server {
	r := a.ChiRouter()

	return http.Server{
		Addr:    addr,
		Handler: r,
	}
}
