// #region place-order-decode
package presentation

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/application"
)

// MakePlaceOrderHandler は POST /orders のハンドラを生成する。
// アプリケーションサービスをクロージャに閉じ込めて返すことで、
// main 関数でのワイヤリングと分離する。
func MakePlaceOrderHandler(
	uc *application.PlaceOrderUseCase,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PlaceOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest,
				"INVALID_JSON", "request body is not valid JSON")
			return
		}
		// #endregion

		// #region place-order-execute
		cmd := application.PlaceOrderCommand{
			CustomerID: req.CustomerID,
			Items:      toApplicationItems(req.Items),
		}

		result, err := uc.Execute(r.Context(), cmd)
		if err != nil {
			writeApplicationError(w, logger, r, err)
			return
		}

		writeJSON(w, http.StatusCreated, PlaceOrderResponse{
			OrderID: result.OrderID,
			Total:   result.Total,
		})
	}
}

// #endregion

// #region get-order-handler
// MakeGetOrderHandler は GET /orders/{id} のハンドラを生成する。
func MakeGetOrderHandler(
	query *application.GetOrderQueryHandler,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id") // Go 1.22+ 標準
		if id == "" {
			writeError(w, http.StatusBadRequest,
				"MISSING_ID", "order id is required")
			return
		}

		result, err := query.Execute(r.Context(),
			application.GetOrderQuery{OrderID: id})
		if err != nil {
			writeApplicationError(w, logger, r, err)
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}

// #endregion

// #region list-orders-handler
// MakeListOrdersHandler は GET /orders のハンドラを生成する。
func MakeListOrdersHandler(
	query *application.ListOrdersByCustomerHandler,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customerID := r.URL.Query().Get("customer_id")
		if customerID == "" {
			writeError(w, http.StatusBadRequest,
				"MISSING_CUSTOMER_ID",
				"customer_id query parameter is required")
			return
		}

		result, err := query.Execute(
			r.Context(),
			application.ListOrdersByCustomerQuery{
				CustomerID: customerID,
			},
		)
		if err != nil {
			writeApplicationError(w, logger, r, err)
			return
		}

		writeJSON(w, http.StatusOK, result)
	}
}

// #endregion
