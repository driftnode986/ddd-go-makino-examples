// #region main-imports
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/infrastructure"
)

// #endregion

// #region main
func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	orderRepo := infrastructure.NewInMemoryOrderRepository()
	pricing := infrastructure.NewStubPricing()
	clock := application.SystemClock{}

	placeOrder := application.NewPlaceOrderUseCase(
		orderRepo, pricing, clock, logger,
	)
	listOrders := &application.ListOrdersByCustomerHandler{Orders: orderRepo}

	_ = placeOrder
	_ = listOrders

	mux := http.NewServeMux()
	// mux.HandleFunc("POST /orders", makePlaceOrderHandler(placeOrder))
	// mux.HandleFunc("GET /orders", makeListOrdersHandler(listOrders))

	srv := &http.Server{Addr: ":8080", Handler: mux}
	logger.InfoContext(ctx, "starting server", slog.String("addr", srv.Addr))
	_ = srv.ListenAndServe()
}

// #endregion
