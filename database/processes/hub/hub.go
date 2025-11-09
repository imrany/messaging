package hub

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, hub *models.Hub) error {
	query := `
		INSERT INTO hubs (name, location, manager_id, description, contact_phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		hub.Name, hub.Location, hub.ManagerID, hub.Description, hub.ContactPhone,
	).Scan(&hub.ID, &hub.CreatedAt)
}

func GetByID(ctx context.Context, id string) (*models.Hub, error) {
	var hub models.Hub
	query := `SELECT id, name, location, manager_id, description, contact_phone, created_at FROM hubs WHERE id = $1`
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&hub.ID, &hub.Name, &hub.Location, &hub.ManagerID, &hub.Description, &hub.ContactPhone, &hub.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("hub not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get hub: %w", err)
	}
	return &hub, nil
}

func Update(ctx context.Context, id string, hub *models.Hub) error {
	query := `
		UPDATE hubs
		SET name = $2, location = $3, manager_id = $4, description = $5, contact_phone = $6
		WHERE id = $1
	`
	_, err := database.DB.ExecContext(ctx, query, id, hub.Name, hub.Location, hub.ManagerID, hub.Description, hub.ContactPhone)
	return err
}

func Delete(ctx context.Context, id string) error {
	query := `DELETE FROM hubs WHERE id = $1`
	_, err := database.DB.ExecContext(ctx, query, id)
	return err
}

func List(ctx context.Context, limit, offset int) ([]models.Hub, error) {
	var hubs []models.Hub
	query := `SELECT id, name, location, manager_id, description, contact_phone, created_at FROM hubs ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := database.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hub models.Hub
		if err := rows.Scan(
			&hub.ID, &hub.Name, &hub.Location, &hub.ManagerID, &hub.Description, &hub.ContactPhone, &hub.CreatedAt,
		); err != nil {
			return nil, err
		}
		hubs = append(hubs, hub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hubs, nil
}
