// #region test-imports
package presentation_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/presentation"
)

// #endregion

// #region place-order-success-request
func TestPlaceOrderHandler_Success(t *testing.T) {
	repo := infrastructure.NewMemoryOrderRepository()
	uc := &application.PlaceOrderUseCase{
		OrderRepo: repo,
		Clock:     application.SystemClock{},
		Logger:    slog.Default(),
	}
	handler := presentation.MakePlaceOrderHandler(uc, slog.Default())

	body := `{
        "customer_id": "cus_001",
        "items": [
            {
                "product_id": "prod_001",
                "quantity": 2,
                "unit_price": 1500
            }
        ]
    }`
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	// #endregion

	// #region place-order-success-assert
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s",
			rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp presentation.PlaceOrderResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.OrderID == "" {
		t.Fatal("order_id is empty")
	}
	if diff := cmp.Diff(int64(3000), resp.Total); diff != "" {
		t.Errorf("total mismatch (-want +got):\n%s", diff)
	}
}

// #endregion

// #region place-order-invalid-json-a
func TestPlaceOrderHandler_InvalidJSON(t *testing.T) {
	repo := infrastructure.NewMemoryOrderRepository()
	uc := &application.PlaceOrderUseCase{
		OrderRepo: repo,
		Clock:     application.SystemClock{},
		Logger:    slog.Default(),
	}
	handler := presentation.MakePlaceOrderHandler(uc, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/orders",
		strings.NewReader(`{invalid json`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}

	// #endregion
	// #region place-order-invalid-json-b
	var resp struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Code != "INVALID_JSON" {
		t.Errorf("code = %q, want %q", resp.Code, "INVALID_JSON")
	}
}

// #endregion

func TestGetOrderHandler_NotFound(t *testing.T) {
	repo := infrastructure.NewMemoryOrderRepository()
	query := &application.GetOrderQueryHandler{OrderRepo: repo}
	handler := presentation.MakeGetOrderHandler(
		query,
		slog.Default(),
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/ord_missing",
		nil,
	)
	req.SetPathValue("id", "ord_missing")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body=%s",
			rec.Code, http.StatusNotFound, rec.Body.String())
	}
}
