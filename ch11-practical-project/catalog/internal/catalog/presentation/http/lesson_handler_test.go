package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/infrastructure/memory"
	httpx "github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/presentation/http"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
)

func TestGetLesson_HTTP(t *testing.T) {
	repo := memory.NewLessonRepository()
	l, err := domain.NewLesson("Pilates", 8, 3000, time.Now(), "trainer-2")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(context.Background(), l); err != nil {
		t.Fatal(err)
	}
	router := httpx.NewRouter(httpx.NewLessonHandler(application.NewLessonQuery(repo)))
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/catalog/lessons/" + l.ID().String())
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var view dto.LessonView
	if err := json.NewDecoder(res.Body).Decode(&view); err != nil {
		t.Fatal(err)
	}
	if view.Name != "Pilates" || view.PriceJPY != 3000 {
		t.Errorf("view = %+v", view)
	}

	missing, err := http.Get(srv.URL + "/catalog/lessons/" + uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}
	_ = missing.Body.Close()
	if missing.StatusCode != http.StatusNotFound {
		t.Errorf("missing lesson: status = %d, want 404", missing.StatusCode)
	}
}
