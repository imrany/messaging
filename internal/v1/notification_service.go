package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/imrany/smart_spore_hub/server/database/models"
	"github.com/imrany/smart_spore_hub/server/database/processes/notification"
)

// GetNotificationPreferences retrieves the notification preferences for a user. - GET /api/v1/notification/preferences/:user_id
func GetNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "user_id is required",
		})
		return
	}

	// Fetch notification preferences from the database
	preferences, err := notification.GetByUserID(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Return the notification preferences as JSON
	if err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Notifcation preference",
		Data:    preferences,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateNotificationPreferences updates the notification preferences for a user. - PUT /api/v1/notification/preferences/:user_id
func UpdateNotificationPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "user_id is required",
		})
		return
	}

	// Parse the request body into a NotificationPreferences struct
	var preferences models.UpdateNotificationPreferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&preferences); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Update the notification preferences in the database
	if err := notification.Update(ctx, userID, preferences); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Return a success response
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Notification preferences updated",
	})
}
