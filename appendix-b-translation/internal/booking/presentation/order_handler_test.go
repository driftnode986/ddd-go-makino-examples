package presentation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	infra "github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

func TestMapErrorToHTTPStatus(t *testing.T) {
	cases := map[error]int{
		domain.ErrOrderNotFound:                         http.StatusNotFound,
		fmt.Errorf("find: %w", domain.ErrOrderNotFound): http.StatusNotFound,
		domain.ErrInvalidStatus:                         http.StatusConflict,
		errors.New("database is down"):                  http.StatusInternalServerError,
	}
	for err, want := range cases {
		if got := mapErrorToHTTPStatus(err); got != want {
			t.Errorf("mapErrorToHTTPStatus(%v) = %d, want %d", err, got, want)
		}
	}
}

func TestOrderHandler(t *testing.T) {
	orders := infra.NewInMemoryOrderRepository()
	uc := application.NewPlaceOrderUseCase(orders, domain.NewPricingService(10),
		test.NewFixedClock(time.Now()))
	srv := httptest.NewServer(NewOrderHandler(uc, orders).Routes())
	t.Cleanup(srv.Close)

	res, err := http.Post(srv.URL+"/orders", "application/json", strings.NewReader(
		`{"customer_id":"cus_001","items":[{"product_id":"p1","quantity":2,"unit_price":1500}]}`))
	if err != nil {
		t.Fatal(err)
	}
	var placed struct {
		OrderID string `json:"order_id"`
		Total   int64  `json:"total"`
	}
	_ = json.NewDecoder(res.Body).Decode(&placed)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusCreated || placed.Total != 3300 {
		t.Fatalf("POST: status = %d, body = %+v", res.StatusCode, placed)
	}

	for path, want := range map[string]int{
		"/orders/" + placed.OrderID: http.StatusOK,
		"/orders/ord_missing":       http.StatusNotFound,
	} {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != want {
			t.Errorf("GET %s: status = %d, want %d", path, res.StatusCode, want)
		}
	}

	bad, err := http.Post(srv.URL+"/orders", "application/json",
		strings.NewReader(`{"customer_id":"cus_001","items":[]}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = bad.Body.Close()
	if bad.StatusCode != http.StatusBadRequest {
		t.Errorf("no items: status = %d, want 400", bad.StatusCode)
	}
}
