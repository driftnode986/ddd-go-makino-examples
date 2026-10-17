package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/booking/internal/booking/domain"
	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	placeUC *application.PlaceBookingUseCase
	listQ   *application.ListBookingsQuery
}

func NewBookingHandler(placeUC *application.PlaceBookingUseCase,
	listQ *application.ListBookingsQuery) *BookingHandler {
	return &BookingHandler{placeUC: placeUC, listQ: listQ}
}

func (h *BookingHandler) Mount(r chi.Router) {
	r.Post("/bookings", h.place)
	r.Get("/members/{member_id}/bookings", h.listByMember)
}

type placeRequest struct {
	MemberID string `json:"member_id"`
	LessonID string `json:"lesson_id"`
}

func (h *BookingHandler) place(w http.ResponseWriter, r *http.Request) {
	var req placeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	out, err := h.placeUC.Execute(r.Context(), application.PlaceBookingInput{
		MemberID: req.MemberID,
		LessonID: req.LessonID,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (h *BookingHandler) listByMember(w http.ResponseWriter, r *http.Request) {
	memberID := chi.URLParam(r, "member_id")
	views, err := h.listQ.ByMember(r.Context(), memberID)
	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, views)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
