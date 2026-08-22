package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticIndexServed(t *testing.T) {
	handler := withStatic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("index status=%d", rec.Code)
	}
}

func TestStaticAppJSServed(t *testing.T) {
	handler := withStatic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/app.js", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("app.js status=%d", rec.Code)
	}
}
