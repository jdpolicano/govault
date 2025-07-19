package tests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/server/middleware"
)

type testBody struct {
	Name string `json:"name"`
}

func TestChainOrder(t *testing.T) {
	order := make([]int, 0, 3)
	a := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, 1)
			next(w, r)
		}
	}
	b := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, 2)
			next(w, r)
		}
	}
	h := func(http.ResponseWriter, *http.Request) { order = append(order, 3) }

	chain := middleware.Chain(h, a, b)
	chain(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("unexpected order %v", order)
	}
}

func TestValidateToken(t *testing.T) {
	ctx := server.NewContext(server.DefaultConfig())
	sess := server.NewSession("bob", []byte("key"), time.Minute)
	ctx.Sessions.Set("id", sess)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer govault-id")

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		called = true
		s, ok := r.Context().Value(server.SessionKey{}).(server.Session)
		if !ok || s.User != "bob" {
			t.Fatalf("session not in context")
		}
	}

	middleware.ValidateToken(ctx)(h)(httptest.NewRecorder(), req)
	if !called {
		t.Fatalf("handler not called")
	}
}

func TestParseJSONBody(t *testing.T) {
	body := bytes.NewBufferString(`{"name":"bob"}`)
	req := httptest.NewRequest("POST", "/", body)

	called := false
	h := func(w http.ResponseWriter, r *http.Request) {
		called = true
		b := r.Context().Value(server.BodyKey{}).(testBody)
		if b.Name != "bob" {
			t.Fatalf("unexpected body: %+v", b)
		}
	}

	middleware.ParseJSONBody[testBody]()(h)(httptest.NewRecorder(), req)
	if !called {
		t.Fatalf("handler not called")
	}
}

func TestLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := log.New(buf, "", 0)

	h := func(http.ResponseWriter, *http.Request) {}
	middleware.Logging(logger)(h)(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	if buf.Len() == 0 {
		t.Fatalf("expected log output")
	}
}
