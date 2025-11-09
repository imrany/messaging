package processes

import (
	"context"
	"fmt"

	"github.com/imrany/gemmie/gemmie-server/pkg/mailer"
	"github.com/imrany/smart_spore_hub/server/database/models"
	"github.com/imrany/smart_spore_hub/server/database/processes/alert"
	"github.com/imrany/smart_spore_hub/server/database/processes/hub"
	"github.com/imrany/smart_spore_hub/server/database/processes/notification"
	"github.com/imrany/smart_spore_hub/server/database/processes/profile"
	sensorreading "github.com/imrany/smart_spore_hub/server/database/processes/sensor_reading"
	"github.com/imrany/smart_spore_hub/server/pkg/whatsapp"
	"github.com/spf13/viper"
)

var SMTP_Config = mailer.SMTPConfig{
	Host:     viper.GetString("SMTP_HOST"),
	Port:     viper.GetInt("SMTP_PORT"),
	Username: viper.GetString("SMTP_USERNAME"),
	Password: viper.GetString("SMTP_PASSWORD"),
	Email:    viper.GetString("SMTP_EMAIL"),
}

// ProcessSensorReading processes a new sensor reading and creates alerts if needed
func ProcessSensorReading(ctx context.Context, req models.CreateSensorReadingRequest) (*models.SensorReading, bool, error) {
	// Insert sensor reading
	reading := &models.SensorReading{
		HubID:       req.HubID,
		Temperature: req.Temperature,
		Humidity:    req.Humidity,
		RecordedAt:  req.RecordedAt,
	}

	if err := sensorreading.Create(ctx, reading); err != nil {
		return nil, false, fmt.Errorf("failed to create sensor reading: %w", err)
	}

	// Check thresholds
	const (
		TEMP_THRESHOLD     = 24.0
		HUMIDITY_THRESHOLD = 65.0
	)

	tempExceeded := req.Temperature > TEMP_THRESHOLD
	humidityExceeded := req.Humidity > HUMIDITY_THRESHOLD

	if !tempExceeded && !humidityExceeded {
		return reading, false, nil
	}

	// Check for existing unresolved alerts
	resolved := false
	existingAlerts, err := alert.List(ctx, models.AlertFilter{
		HubID:    req.HubID,
		Resolved: resolved,
		Limit:    1,
	})

	if err != nil {
		return reading, false, fmt.Errorf("failed to check existing alerts: %w", err)
	}

	// Only create alert if no unresolved alert exists
	if len(existingAlerts) == 0 {
		var alertType models.AlertType
		var message string

		if tempExceeded && humidityExceeded {
			alertType = models.AlertTypeBoth
			message = fmt.Sprintf("ALERT: Both temperature (%.2f°C) and humidity (%.2f%%) have exceeded safe thresholds!", req.Temperature, req.Humidity)
		} else if tempExceeded {
			alertType = models.AlertTypeTemperature
			message = fmt.Sprintf("ALERT: Temperature (%.2f°C) has exceeded the safe threshold of %.2f°C!", req.Temperature, TEMP_THRESHOLD)
		} else {
			alertType = models.AlertTypeHumidity
			message = fmt.Sprintf("ALERT: Humidity (%.2f%%) has exceeded the safe threshold of %.2f%%!", req.Humidity, HUMIDITY_THRESHOLD)
		}

		alertMsg := &models.Alert{
			HubID:       req.HubID,
			AlertType:   alertType,
			Message:     message,
			Temperature: req.Temperature,
			Humidity:    req.Humidity,
		}

		if err := alert.Create(ctx, alertMsg); err != nil {
			return reading, true, fmt.Errorf("failed to create alert: %w", err)
		}

		// Send notifications
		go sendAlertNotifications(context.Background(), alertMsg)
	}

	return reading, tempExceeded || humidityExceeded, nil
}

// GetSensorData retrieves the latest sensor data for a given hub
func GetSensorData(ctx context.Context, hubID string) ([]*models.SensorReading, error) {
	readings, err := sensorreading.GetLatestByHub(ctx, hubID)
	if err != nil {
		return nil, err
	}

	return []*models.SensorReading{readings}, nil
}

func GetUnresolvedAlerts(ctx context.Context, hubID string) ([]models.Alert, error) {
	alerts, err := alert.GetUnresolvedByHub(ctx, hubID)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

// sendAlertNotifications sends notifications to the hub manager
func sendAlertNotifications(ctx context.Context, alert *models.Alert) error {
	// Get hub details
	hub, err := hub.GetByID(ctx, alert.HubID)
	if err != nil || hub.ManagerID == "" {
		return fmt.Errorf("%s", err)
	}

	// Get notification preferences
	prefs, err := notification.GetByUserID(ctx, hub.ManagerID)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	// Get profile details
	profile, err := profile.GetByID(ctx, hub.ManagerID)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	// Send notifications based on preferences
	if prefs.WhatsAppEnabled && prefs.PhoneNumber != "" {
		// TODO: Send WhatsApp message
		whatsapp.SendMessage(ctx, prefs.PhoneNumber, alert.Message)
	}

	if prefs.SMSEnabled && prefs.PhoneNumber != "" {
		// TODO: Send SMS
		// sms.Send(ctx, *prefs.PhoneNumber, alert.Message)
	}

	if prefs.EmailEnabled {
		// TODO: Send email
		emailData := mailer.EmailData{
			To: []string{
				profile.Email,
			},
			Subject: "Alert: " + hub.Name,
			Body:    alert.Message,
			IsHTML:  false,
		}
		mailer.SendEmail(emailData, SMTP_Config)
	}

	return nil
}
