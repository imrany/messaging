package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/imrany/smart_spore_hub/server/database/processes/hub"
)

// GetUserHubs retrieves all hubs for a user - GET /v1/hubs/{user_id}
func GetUserHubs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	userID := chi.URLParam(r, "user_id")
	w.Header().Set("Content-Type", "application/json")
	hubs, err := hub.GetByID(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Message: "Hubs retrieved successfully",
		Data:    hubs,
	})
}

// GetHubs retrieves all hubs - GET /v1/hubs
func GetHubs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	hubs, err := hub.List(ctx, 50, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Message: "Hubs retrieved successfully",
		Data:    hubs,
	})
}
