// catalog/cmd/api/main.go
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

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/infrastructure/memory"
	httpx "github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/presentation/http"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	repo := memory.NewLessonRepository()
	query := application.NewLessonQuery(repo)
	handler := httpx.NewLessonHandler(query)
	router := httpx.NewRouter(handler)

	srv := &http.Server{
		Addr:              ":8082",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("catalog BC listening", "addr", srv.Addr)
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
