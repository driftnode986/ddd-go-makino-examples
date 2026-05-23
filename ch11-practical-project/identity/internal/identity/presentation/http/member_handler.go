// identity/internal/identity/presentation/http/member_handler.go
package http

import (
	"encoding/json"
	"net/http"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/application"
	"github.com/go-chi/chi/v5"
)

type MemberHandler struct {
	registerUC *application.RegisterMemberUseCase
}

func NewMemberHandler(uc *application.RegisterMemberUseCase) *MemberHandler {
	return &MemberHandler{registerUC: uc}
}

func (h *MemberHandler) Mount(r chi.Router) {
	r.Post("/members", h.register)
}

func (h *MemberHandler) register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	out, err := h.registerUC.Execute(r.Context(), application.RegisterMemberInput{
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}

func NewRouter(h *MemberHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/identity", h.Mount)
	return r
}
