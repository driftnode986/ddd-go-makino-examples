package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

// bookLessonExecutor は handler が使うユースケースの契約 (テストで差し替える)。
type bookLessonExecutor interface {
	Execute(
		ctx context.Context,
		in application.BookLessonInput,
	) error
}

// Handler は予約 API の HTTP ハンドラ。
type Handler struct {
	useCase bookLessonExecutor
}

// NewHandler は *application.BookLessonUseCase などを受け取ってハンドラを作る。
func NewHandler(useCase bookLessonExecutor) *Handler {
	return &Handler{useCase: useCase}
}

// ServeHTTP は POST /bookings を handleCreateBooking に渡す。
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}
	h.handleCreateBooking(w, r)
}

// buildInput はリクエストの JSON (CreateBookingRequest) をユースケースの入力に変換する。
// 読めない値はゼロ値のまま渡し、ドメイン層の NewBooking に拒ませる。
func buildInput(r *http.Request) application.BookLessonInput {
	var req CreateBookingRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	at, _ := time.Parse(time.RFC3339, req.ScheduledAt)
	return application.BookLessonInput{
		MemberID:    domain.MemberID(req.MemberID),
		LessonID:    domain.LessonID(req.LessonID),
		ScheduledAt: at,
	}
}
