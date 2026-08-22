package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPWorkflow(t *testing.T) {
	ts := httptest.NewServer(NewRouter(NewCycleService(NewCellStore())))
	defer ts.Close()
	for _, path := range []string{"/health", "/api/cells", "/"} {
		resp, err := http.Get(ts.URL + path)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("GET %s status=%v err=%v", path, resp.StatusCode, err)
		}
		resp.Body.Close()
	}
	cases := []struct {
		body string
		want int
	}{{`{"id":"cell-a14","status":"complete"}`, http.StatusOK}, {`{"id":"cell-a14","status":"failed"}`, http.StatusBadRequest}, {`{"id":"cell-none","status":"paused"}`, http.StatusNotFound}, {`not-json`, http.StatusBadRequest}}
	for _, tc := range cases {
		resp, err := http.Post(ts.URL+"/api/cells/status", "application/json", bytes.NewBufferString(tc.body))
		if err != nil || resp.StatusCode != tc.want {
			t.Fatalf("POST %s status=%v want=%v err=%v", tc.body, resp.StatusCode, tc.want, err)
		}
		resp.Body.Close()
	}
}

func postStatus(t *testing.T, body string) int {
	t.Helper()
	ts := httptest.NewServer(NewRouter(NewCycleService(NewCellStore())))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/cells/status", "application/json", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func TestHTTPUnknownCellReturnsNotFound(t *testing.T) {
	if got := postStatus(t, `{"id":"cell-none","status":"paused"}`); got != http.StatusNotFound {
		t.Fatalf("status=%d want 404", got)
	}
}

func TestHTTPInvalidStatusRejected400(t *testing.T) {
	if got := postStatus(t, `{"id":"cell-a14","status":"failed"}`); got != http.StatusBadRequest {
		t.Fatalf("status=%d want 400", got)
	}
}

func TestHTTPInvalidBodyReturnsBadRequest(t *testing.T) {
	if got := postStatus(t, `not-json`); got != http.StatusBadRequest {
		t.Fatalf("status=%d want 400", got)
	}
}
