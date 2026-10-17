package presentation_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/presentation"
)

// placeOrder は POST /orders を通して注文を 1 件作り、注文 ID を返す。
func placeOrder(
	t *testing.T,
	repo *infrastructure.MemoryOrderRepository,
) string {
	t.Helper()
	uc := &application.PlaceOrderUseCase{
		OrderRepo: repo,
		Clock:     application.SystemClock{},
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	body := `{"customer_id":"cus_001","items":[` +
		`{"product_id":"p1","quantity":1,"unit_price":800}]}`
	rec := httptest.NewRecorder()
	presentation.MakePlaceOrderHandler(uc, slog.Default()).
		ServeHTTP(
			rec,
			httptest.NewRequest(
				http.MethodPost,
				"/orders",
				strings.NewReader(body),
			),
		)
	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"place order: status = %d, body=%s",
			rec.Code,
			rec.Body.String(),
		)
	}
	var resp presentation.PlaceOrderResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	return resp.OrderID
}

func TestGetOrderHandler_Found(t *testing.T) {
	repo := infrastructure.NewMemoryOrderRepository()
	id := placeOrder(t, repo)
	handler := presentation.MakeGetOrderHandler(
		&application.GetOrderQueryHandler{
			OrderRepo: repo,
		},
		slog.Default(),
	)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+id, nil)
	req.SetPathValue("id", id)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var got application.OrderSummary
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.OrderID != id || got.Total != 800 {
		t.Errorf("summary = %+v", got)
	}
}

func TestListOrdersHandler(t *testing.T) {
	repo := infrastructure.NewMemoryOrderRepository()
	placeOrder(t, repo)
	handler := presentation.MakeListOrdersHandler(
		&application.ListOrdersByCustomerHandler{
			OrderRepo: repo,
		},
		slog.Default(),
	)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(
		rec,
		httptest.NewRequest(
			http.MethodGet,
			"/orders?customer_id=cus_001",
			nil,
		),
	)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var got []application.OrderSummary
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Errorf("orders = %d, want 1", len(got))
	}

	rec = httptest.NewRecorder()
	handler.ServeHTTP(
		rec,
		httptest.NewRequest(http.MethodGet, "/orders", nil),
	)
	if rec.Code != http.StatusBadRequest {
		t.Errorf(
			"without customer_id: status = %d, want 400",
			rec.Code,
		)
	}
}

func TestLoggingMiddleware_RecordsStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		},
	)
	handler := presentation.RequestIDMiddleware(
		presentation.LoggingMiddleware(logger)(next),
	)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-Request-ID", "req-1")
	handler.ServeHTTP(httptest.NewRecorder(), req)

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log is not JSON: %q", buf.String())
	}
	if entry["status"] != float64(http.StatusTeapot) ||
		entry["request_id"] != "req-1" {
		t.Errorf("log entry = %v", entry)
	}
}
