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

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/presentation"
)

// #endregion

// #region main-wiring
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo}))

	// Infrastructure 層: in-memory リポジトリ
	orderRepo := infrastructure.NewMemoryOrderRepository()

	// Application 層: ユースケースとクエリ
	placeOrderUC := &application.PlaceOrderUseCase{
		OrderRepo: orderRepo,
		Clock:     application.SystemClock{},
		Logger:    logger,
	}
	getOrderQuery := &application.GetOrderQueryHandler{
		OrderRepo: orderRepo,
	}
	listOrdersQuery := &application.ListOrdersByCustomerHandler{
		OrderRepo: orderRepo,
	}
	// #endregion

	// #region main-routes
	// Presentation 層: ハンドラ
	// #region routes
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders",
		presentation.MakePlaceOrderHandler(placeOrderUC, logger))
	mux.HandleFunc("GET /orders/{id}",
		presentation.MakeGetOrderHandler(getOrderQuery, logger))
	mux.HandleFunc("GET /orders",
		presentation.MakeListOrdersHandler(listOrdersQuery, logger))
	// #endregion

	// ミドルウェア合成
	handler := presentation.RequestIDMiddleware(
		presentation.LoggingMiddleware(logger)(
			presentation.RecoveryMiddleware(logger)(mux),
		),
	)
	// #endregion

	// #region server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// 起動 + graceful shutdown
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// #endregion

	// #region shutdown
	select {
	case sig := <-sigCh:
		logger.Info(
			"shutdown signal received",
			slog.String("signal", sig.String()),
		)
	case err := <-serverErr:
		logger.Error("server failed", slog.Any("err", err))
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", slog.Any("err", err))
	}
}

// #endregion
