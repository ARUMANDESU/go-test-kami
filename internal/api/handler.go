package api

import (
	"errors"
	"net/http"

	"github.com/ARUMANDESU/go-test-kami/internal/domain"
	"github.com/ARUMANDESU/go-test-kami/internal/service/reservation"
	"github.com/go-chi/chi/v5"
)

// NotFound is custom 404 handler
func (api API) NotFound(w http.ResponseWriter, r *http.Request) {
	api.notFoundResponse(w, r)
}

// Healthcheck returns the status of the service
func (api API) Healthcheck(w http.ResponseWriter, r *http.Request) {
	api.writeJSON(w, http.StatusOK, envelope{"status": "available"}, nil)
}

// ReserveRoom reserves a room.
//
//	POST /v1/reservations
//
//	{
//		"room_id": "string",
//		"start_time": "string", // RFC3339
//		"end_time": "string" // RFC3339
//	}
func (api API) ReserveRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var dto domain.ReservationCreateDTO

	api.readJSON(w, r, &dto)

	reservedRoom, err := api.ReservationService.ReserveRoom(ctx, dto)
	if err != nil {
		switch {
		case errors.Is(err, reservation.ErrReservationConflict):
			api.conflictResponse(w, r)
		case errors.Is(err, reservation.ErrInvalidArgument):
			api.failedValidationResponse(w, r, err)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}

	api.writeJSON(w, http.StatusCreated, envelope{"reservation": reservedRoom}, nil)
}

// GetRoomReservations returns all reservations for a room.
//
//	GET /v1/reservations/{room_id}
func (api API) GetRoomReservations(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	if roomID == "" {
		api.failedValidationResponse(w, r, errors.New("missing room_id parameter"))
		return
	}

	ctx := r.Context()
	reservations, err := api.ReservationService.GetRoomReservations(ctx, roomID)
	if err != nil {
		switch {
		case errors.Is(err, reservation.ErrInvalidArgument):
			api.failedValidationResponse(w, r, err)
		default:
			api.serverErrorResponse(w, r, err)
		}
		return
	}

	api.writeJSON(w, http.StatusOK, envelope{"reservations": reservations}, nil)
}
