// #region main-imports
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/infrastructure/memory"
	httpx "github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/presentation/http"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
)

type stubLessonQuery struct{}

func (stubLessonQuery) FindLessonView(_ context.Context,
	id string) (dto.LessonView, error) {
	return dto.LessonView{ID: id, Name: "Sample Lesson"}, nil
}

// #endregion

// #region main-wiring
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := memory.NewBookingRepository()
	placeUC := application.NewPlaceBookingUseCase(repo, stubLessonQuery{}, time.Now)
	listQ := application.NewListBookingsQuery(repo)
	handler := httpx.NewBookingHandler(placeUC, listQ)
	router := httpx.NewRouter(handler, logger)

	srv := &http.Server{
		Addr:              ":8081",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	// #endregion

	// #region main-shutdown
	go func() {
		logger.Info("booking BC listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", "err", err)
	}
}

// #endregion
