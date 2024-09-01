package api

import (
	"errors"
	"net/http"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
	"github.com/ARUMANDESU/go-test-kami/internal/service/reservation"
	"github.com/segmentio/encoding/json"
)

func (a API) ReserveRoom(w http.ResponseWriter, r *http.Request) {

	var dto domain.ReservationCreateDTO
	// decode request body to dto
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	reservedRoom, err := a.ReservationService.ReserveRoom(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, reservation.ErrReservationConflict):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, reservation.ErrInvalidArgument):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// encode response
	if err := json.NewEncoder(w).Encode(reservedRoom); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a API) GetRoomReservations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	roomID := r.URL.Query().Get("room_id")
	reservations, err := a.ReservationService.GetRoomReservations(ctx, roomID)
	if err != nil {
		switch {
		case errors.Is(err, reservation.ErrInvalidArgument):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(reservations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
