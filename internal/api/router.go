package api

import (
	"net/http"

	"user-restriction-manager/internal/api/handler"
	"user-restriction-manager/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	restrictionHandler *handler.RestrictionHandler,
	appealHandler *handler.AppealHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/restrictions", func(r chi.Router) {
			r.Post("/", restrictionHandler.Create)
			r.Get("/{userID}", restrictionHandler.GetActive)
			r.Delete("/{id}", restrictionHandler.Revoke)
			r.Post("/{restrictionID}/appeal", appealHandler.Submit) // Alternate path, but sticking to plan
		})

		r.Route("/appeals", func(r chi.Router) {
			r.Post("/", appealHandler.Submit)
			r.Post("/{id}/review", appealHandler.Review)
		})
	})

	return r
}
