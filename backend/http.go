package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type statusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewRouter(service *CycleService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/cells", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, service.Cells()) })
	mux.HandleFunc("POST /api/cells/status", func(w http.ResponseWriter, r *http.Request) {
		var req statusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id and status JSON are required"})
			return
		}
		if strings.TrimSpace(req.ID) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id and status JSON are required"})
			return
		}
		cell, err := service.ChangeStatus(req.ID, req.Status)
		if err != nil {
			writeJSON(w, statusForChangeStatusError(err), map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, cell)
	})
	return withStatic(mux)
}

// statusForChangeStatusError maps a ChangeStatus error to the matching HTTP
// status code. An unknown cell is a client error (404), an invalid status is a
// client error (400); anything else falls back to 500.
func statusForChangeStatusError(err error) int {
	if errors.Is(err, ErrCellNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrCellStatusInvalid) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
