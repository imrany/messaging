package v1

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler handles the health check endpoint. GET /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	version := "v1.0.0"

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Server is healthy",
		Data: map[string]any{
			"timestamp":   time.Now(),
			"version":     version,
			"description": "A smart spore hub mail server",
		},
	})
}
