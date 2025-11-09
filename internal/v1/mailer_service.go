package v1

import (
	"encoding/json"
	"net/http"

	"github.com/imrany/gemmie/gemmie-server/pkg/mailer"
	"github.com/spf13/viper"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

type SendMailRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	IsHTML  bool     `json:"is_html"`
}

var SMTP_Config = mailer.SMTPConfig{
	Host:     viper.GetString("SMTP_HOST"),
	Port:     viper.GetInt("SMTP_PORT"),
	Username: viper.GetString("SMTP_USERNAME"),
	Password: viper.GetString("SMTP_PASSWORD"),
	Email:    viper.GetString("SMTP_EMAIL"),
}

// SendMail handler - POST /api/v1/mailer/send - receives send email requests
func SendMail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SendMailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	emailData := mailer.EmailData{
		To:      req.To,
		Subject: "Test Email",
		Body:    "This is a test email.",
		IsHTML:  false,
	}
	err := mailer.SendEmail(emailData, SMTP_Config)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to send email",
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Email sent successfully",
	})
}
