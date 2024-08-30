package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

func (a *API) ChiRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(slogchi.New(a.log))
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

		v1.Post("/reservation", a.ReserveRoom)

		v1.Get("/reservation/{roomID}", a.GetRoomReservations)
	})

	return r
}
