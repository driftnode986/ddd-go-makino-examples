package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
	"github.com/go-chi/chi/v5"
)

type LessonHandler struct {
	query *application.LessonQuery
}

func NewLessonHandler(query *application.LessonQuery) *LessonHandler {
	return &LessonHandler{query: query}
}

func (h *LessonHandler) Mount(r chi.Router) {
	r.Get("/lessons/{id}", h.get)
}

func (h *LessonHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	view, err := h.query.FindLessonView(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrLessonNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(view)
}

func NewRouter(handler *LessonHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/catalog", handler.Mount)
	return r
}
