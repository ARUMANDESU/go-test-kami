package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

func (a *API) ChiRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(slogchi.New(a.log))
	r.Use(middleware.Recoverer)

	r.NotFound(a.NotFound)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", a.Healthcheck)
		r.Get("/reservations/{room_id}", a.GetRoomReservations)

		r.Post("/reservations", a.ReserveRoom)
	})

	return r
}
