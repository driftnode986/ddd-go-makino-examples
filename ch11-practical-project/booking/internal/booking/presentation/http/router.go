package http

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	handler *BookingHandler,
	logger *slog.Logger,
) chi.Router {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(loggingMiddleware(logger))
	r.Use(chimiddleware.Recoverer) // ロギングより内側（第9章）

	r.Route("/booking", handler.Mount)
	return r
}
