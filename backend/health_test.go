package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthStatusOK(t *testing.T) {
	rec := httptest.NewRecorder()
	healthHandler(rec, httptest.NewRequest("GET", "/health", nil))
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if rec.Code != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("status=%d body=%v", rec.Code, body)
	}
}

func TestHealthServiceName(t *testing.T) {
	rec := httptest.NewRecorder()
	healthHandler(rec, httptest.NewRequest("GET", "/health", nil))
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["service"] != "battery-lab-cycle" {
		t.Fatalf("service=%q", body["service"])
	}
}

func TestHealthVersion(t *testing.T) {
	rec := httptest.NewRecorder()
	healthHandler(rec, httptest.NewRequest("GET", "/health", nil))
	var body map[string]string
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["version"] != "1.0.0" {
		t.Fatalf("version=%q", body["version"])
	}
}
