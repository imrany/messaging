package alert

import (
	"context"
	"fmt"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, alert *models.Alert) error {
	query := `
		INSERT INTO alerts (hub_id, alert_type, message, temperature, humidity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		alert.HubID, alert.AlertType, alert.Message, alert.Temperature, alert.Humidity,
	).Scan(&alert.ID, &alert.CreatedAt)
}

func GetByID(ctx context.Context, id string) (*models.Alert, error) {
	var alert models.Alert
	query := `SELECT hub_id, alert_type, message, temperature, humidity, created_at, resolved, resolved_at, id FROM alerts WHERE id = $1`
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&alert.HubID, &alert.AlertType, &alert.Message, &alert.Temperature, &alert.Humidity, &alert.CreatedAt, &alert.Resolved, &alert.ResolvedAt, &alert.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("alert not found")
	}
	return &alert, err
}

func Resolve(ctx context.Context, id string) error {
	query := `UPDATE alerts SET resolved = true, resolved_at = NOW() WHERE id = $1`
	_, err := database.DB.ExecContext(ctx, query, id)
	return err
}

func List(ctx context.Context, filter models.AlertFilter) ([]models.Alert, error) {
	var alerts []models.Alert
	query := `
		SELECT id, hub_id, alert_type, message, temperature, humidity, created_at, resolved, resolved_at FROM alerts
		WHERE (hub_id = $1 OR $1::uuid IS NULL)
				AND (alert_type = $2 OR $2::text IS NULL)
				AND (resolved = $3 OR $3::boolean IS NULL)
				AND (created_at >= $4 OR $4::timestamptz IS NULL)
				AND (created_at <= $5 OR $5::timestamptz IS NULL)
		ORDER BY created_at DESC
		LIMIT $6 OFFSET $7
	`

	rows, err := database.DB.QueryContext(ctx, query, filter.HubID, filter.AlertType, filter.Resolved, filter.StartDate, filter.EndDate, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alert models.Alert
		if err := rows.Scan(&alert.ID, &alert.HubID, &alert.AlertType, &alert.Message, &alert.Temperature, &alert.Humidity, &alert.CreatedAt, &alert.Resolved, &alert.ResolvedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func GetUnresolvedByHub(ctx context.Context, hubID string) ([]models.Alert, error) {
	var alerts []models.Alert
	query := `SELECT id, hub_id, alert_type, message, temperature, humidity, created_at, resolved, resolved_at FROM alerts WHERE hub_id = $1 AND resolved = false ORDER BY created_at DESC`

	rows, err := database.DB.QueryContext(ctx, query, hubID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alert models.Alert
		if err := rows.Scan(&alert.ID, &alert.HubID, &alert.AlertType, &alert.Message, &alert.Temperature, &alert.Humidity, &alert.CreatedAt, &alert.Resolved, &alert.ResolvedAt); err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}
