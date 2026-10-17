package presentation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
)

// OrderHandler は注文の HTTP API（POST /orders、GET /orders/{id}）。
type OrderHandler struct {
	placeOrder *application.PlaceOrderUseCase
	orders     domain.OrderRepository
}

func NewOrderHandler(placeOrder *application.PlaceOrderUseCase, orders domain.OrderRepository) *OrderHandler {
	return &OrderHandler{placeOrder: placeOrder, orders: orders}
}

func (h *OrderHandler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.place)
	mux.HandleFunc("GET /orders/{id}", h.get)
	return mux
}

type placeOrderRequest struct {
	CustomerID string `json:"customer_id"`
	Items      []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
		UnitPrice int64  `json:"unit_price"`
	} `json:"items"`
}

func (h *OrderHandler) place(w http.ResponseWriter, r *http.Request) {
	var req placeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	cmd := application.PlaceOrderCommand{CustomerID: domain.CustomerID(req.CustomerID)}
	for _, it := range req.Items {
		cmd.Items = append(cmd.Items, domain.OrderItem{
			ProductID: it.ProductID, Quantity: it.Quantity, UnitPrice: it.UnitPrice,
		})
	}
	res, err := h.placeOrder.Execute(r.Context(), cmd)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"order_id": res.OrderID, "total": res.Total})
}

func (h *OrderHandler) get(w http.ResponseWriter, r *http.Request) {
	o, err := h.orders.FindByID(r.Context(), domain.OrderID(r.PathValue("id")))
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"order_id": o.ID(), "status": o.Status()})
}

// writeError は入力の誤りを 400 にし、それ以外は mapErrorToHTTPStatus に任せる。
func (h *OrderHandler) writeError(w http.ResponseWriter, err error) {
	status := mapErrorToHTTPStatus(err)
	if errors.Is(err, domain.ErrInvalidOrder) {
		status = http.StatusBadRequest
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
