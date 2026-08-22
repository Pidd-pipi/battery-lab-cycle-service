package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnterpriseServerHasRequestTimeouts(t *testing.T) {
	srv := newEnterpriseServer(":0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if srv.ReadHeaderTimeout <= 0 {
		t.Fatal("ReadHeaderTimeout not set")
	}
	if srv.ReadTimeout <= 0 {
		t.Fatal("ReadTimeout not set")
	}
	if srv.WriteTimeout <= 0 {
		t.Fatal("WriteTimeout not set")
	}
	if srv.IdleTimeout <= 0 {
		t.Fatal("IdleTimeout not set")
	}
}

func TestRequestIDHeadersUnique(t *testing.T) {
	handler := requestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	seen := map[string]bool{}
	for i := 0; i < 5; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(rec, req)
		id := rec.Header().Get("X-Request-ID")
		if id == "" {
			t.Fatal("missing request id")
		}
		if seen[id] {
			t.Fatalf("duplicate request id %q", id)
		}
		seen[id] = true
	}
}
