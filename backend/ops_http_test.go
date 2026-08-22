package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpsMiddlewareSetsLatencyHeader(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/api/cells", nil))
	if rec.Header().Get("X-Operations-Latency-Ms") == "" {
		t.Fatal("missing latency header")
	}
}

func TestOpsMiddlewareSetsDomainHeader(t *testing.T) {
	handler := opsEnterpriseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/api/cells", nil))
	if got := rec.Header().Get("X-Operations-Domain"); got != opsDomainName {
		t.Fatalf("domain=%q want %q", got, opsDomainName)
	}
}

func TestOpsJSONSetsContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	opsJSON(rec, http.StatusOK, map[string]string{"ok": "1"})
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type=%q", ct)
	}
}

func TestOpsNoStoreHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	opsNoStore(rec)
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("cache-control=%q", rec.Header().Get("Cache-Control"))
	}
}
