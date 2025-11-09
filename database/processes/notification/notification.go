package notification

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, pref *models.NotificationPreference) error {
	query := `
		INSERT INTO notification_preferences (user_id, sms_enabled, whatsapp_enabled, email_enabled, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		pref.UserID, pref.SMSEnabled, pref.WhatsAppEnabled, pref.EmailEnabled, pref.PhoneNumber,
	).Scan(&pref.ID, &pref.CreatedAt)
}

func GetByUserID(ctx context.Context, userID string) (*models.NotificationPreference, error) {
	var pref models.NotificationPreference
	query := `SELECT id, user_id, sms_enabled, whatsapp_enabled, email_enabled, phone_number, created_at, email FROM notification_preferences WHERE user_id = $1`
	err := database.DB.QueryRowContext(ctx, query, userID).Scan(
		&pref.ID,
		&pref.UserID,
		&pref.SMSEnabled,
		&pref.WhatsAppEnabled,
		&pref.EmailEnabled,
		&pref.PhoneNumber,
		&pref.CreatedAt,
		&pref.Email,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("notification preferences not found")
	}
	return &pref, err
}

func Update(ctx context.Context, userID string, req models.UpdateNotificationPreferenceRequest) error {
	query := `
		UPDATE notification_preferences
		SET sms_enabled = COALESCE($2, sms_enabled),
		    whatsapp_enabled = COALESCE($3, whatsapp_enabled),
		    email_enabled = COALESCE($4, email_enabled),
		    phone_number = COALESCE($5, phone_number)
		WHERE user_id = $1
	`
	_, err := database.DB.ExecContext(ctx, query, userID, req.SMSEnabled, req.WhatsAppEnabled, req.EmailEnabled, req.PhoneNumber)
	return err
}

func Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM notification_preferences WHERE user_id = $1`
	_, err := database.DB.ExecContext(ctx, query, userID)
	return err
}
