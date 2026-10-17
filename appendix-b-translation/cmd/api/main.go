package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	infra "github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/presentation"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/system"
)

// #region main
func main() {
	db := mustOpenDB()
	clock := system.NewClock()
	orders := infra.NewOrderRepositoryPG(db)
	pricing := domain.NewPricingService(10)
	placeOrder := application.NewPlaceOrderUseCase(
		orders,
		pricing,
		clock,
	)
	// ...HTTP ハンドラの組み立てに続く
	// #endregion
	handler := presentation.NewOrderHandler(placeOrder, orders)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

// mustOpenDB は環境変数 DATABASE_URL の PostgreSQL に接続する。
// テーブルは db/booking_schema.sql で作っておく。
func mustOpenDB() *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal(
			"DATABASE_URL is not set " +
				"(例: postgres://user:pass@localhost:5432/db)",
		)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	return pool
}
