package v1

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler returns server health status - GET /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status      string `json:"status"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Uptime      string `json:"uptime"`
	}{
		Status:      http.StatusText(r.Response.StatusCode),
		Version:     "1.0.0",
		Description: "Service is healthy",
		Uptime:      time.Since(time.Now()).String(),
	})
}
