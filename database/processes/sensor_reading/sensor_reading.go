package sensorreading

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, reading *models.SensorReading) error {
	query := `
		INSERT INTO sensor_readings (hub_id, temperature, humidity, recorded_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		reading.HubID, reading.Temperature, reading.Humidity, reading.RecordedAt,
	).Scan(&reading.ID, &reading.CreatedAt)
}

func GetByID(ctx context.Context, id string) (*models.SensorReading, error) {
	var reading models.SensorReading
	query := `SELECT id, hub_id, temperature, humidity, recorded_at, created_at FROM sensor_readings WHERE id = $1`
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&reading.ID, &reading.HubID, &reading.Temperature, &reading.Humidity, &reading.RecordedAt, &reading.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sensor reading not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get sensor reading: %w", err)
	}
	return &reading, nil
}

func List(ctx context.Context, filter models.SensorReadingFilter) ([]models.SensorReading, error) {
	var readings []models.SensorReading
	query := `
		SELECT id, hub_id, temperature, humidity, recorded_at, created_at FROM sensor_readings
		WHERE ($1::uuid IS NULL OR hub_id = $1)
				AND ($2::timestamptz IS NULL OR recorded_at >= $2)
				AND ($3::timestamptz IS NULL OR recorded_at <= $3)
		ORDER BY recorded_at DESC
		LIMIT $4 OFFSET $5
	`
	rows, err := database.DB.QueryContext(ctx, query, filter.HubID, filter.StartDate, filter.EndDate, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sensor readings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var reading models.SensorReading
		if err := rows.Scan(
			&reading.ID, &reading.HubID, &reading.Temperature, &reading.Humidity, &reading.RecordedAt, &reading.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sensor reading: %w", err)
		}
		readings = append(readings, reading)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sensor readings: %w", err)
	}
	return readings, nil
}

func GetLatestByHub(ctx context.Context, hubID string) (*models.SensorReading, error) {
	var reading models.SensorReading
	query := `SELECT id, hub_id, temperature, humidity, recorded_at, created_at FROM sensor_readings WHERE hub_id = $1 ORDER BY recorded_at DESC LIMIT 1`
	err := database.DB.QueryRowContext(ctx, query, hubID).Scan(
		&reading.ID, &reading.HubID, &reading.Temperature, &reading.Humidity, &reading.RecordedAt, &reading.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no sensor readings found for hub")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest sensor reading: %w", err)
	}
	return &reading, nil
}
