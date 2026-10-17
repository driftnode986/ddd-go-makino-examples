package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/infrastructure/memory"
	httpx "github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/presentation/http"
)

func TestRegisterMember_HTTP(t *testing.T) {
	uc := application.NewRegisterMemberUseCase(memory.NewMemberRepository())
	srv := httptest.NewServer(httpx.NewRouter(httpx.NewMemberHandler(uc)))
	t.Cleanup(srv.Close)

	res, err := http.Post(srv.URL+"/identity/members", "application/json",
		strings.NewReader(`{"email":"jiro@example.com","name":"Jiro"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want 201", res.StatusCode)
	}

	bad, err := http.Post(srv.URL+"/identity/members", "application/json",
		strings.NewReader(`{"email":"broken","name":"Jiro"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = bad.Body.Close()
	if bad.StatusCode != http.StatusBadRequest {
		t.Errorf("invalid email: status = %d, want 400", bad.StatusCode)
	}
}
