package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/imrany/whats-email/pkg/whatsapp"
)

type WhatsAppRequest struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// SendWhatsAppMessage sends a WhatsApp message using the WhatsApp API - POST /api/v1/whatsapp/send.
func SendWhatsAppMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req WhatsAppRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	err = whatsapp.SendMessage(r.Context(), req.PhoneNumber, req.Message)
	if err != nil {
		slog.Error("Failed to send WhatsApp message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to send WhatsApp message",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "WhatsApp message sent successfully",
	})
}
