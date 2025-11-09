package profile

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, profile *models.Profile) error {
	query := `
		INSERT INTO profiles (id, full_name, phone, role, location, created_at, updated_at, email, password)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		profile.ID, profile.FullName, profile.Phone, profile.Role,
		profile.Location, time.Now(), time.Now(), profile.Email, profile.Password,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
}

func GetByID(ctx context.Context, id string) (*models.Profile, error) {
	var profile models.Profile
	query := `SELECT id, full_name, phone, role, location, created_at, updated_at, email, password FROM profiles WHERE id = $1`
	row := database.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&profile.ID, &profile.FullName, &profile.Phone, &profile.Role, &profile.Location, &profile.CreatedAt, &profile.UpdatedAt, &profile.Email, &profile.Password)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile not found")
	} else if err != nil {
		return nil, err
	}

	return &profile, nil
}

func GetByEmail(ctx context.Context, email string) (*models.Profile, error) {
	var profile models.Profile
	query := `SELECT id, full_name, phone, role, location, created_at, updated_at, email, password FROM profiles WHERE email = $1`
	row := database.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&profile.ID, &profile.FullName, &profile.Phone, &profile.Role, &profile.Location, &profile.CreatedAt, &profile.UpdatedAt, &profile.Email, &profile.Password)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("profile with email %s not found", email)
	} else if err != nil {
		return nil, err
	}

	return &profile, nil
}

func Update(ctx context.Context, id string, req models.UpdateProfileRequest) error {
	query := `
		UPDATE profiles
		SET full_name = COALESCE($2, full_name),
		    phone = COALESCE($3, phone),
		    role = COALESCE($4, role),
		    location = COALESCE($5, location),
			password = COALESCE($6, password),
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := database.DB.ExecContext(ctx, query, id, req.FullName, req.Phone, req.Role, req.Location, req.Password)
	return err
}

func Delete(ctx context.Context, id string) error {
	query := `DELETE FROM profiles WHERE id = $1`
	_, err := database.DB.ExecContext(ctx, query, id)
	return err
}

func List(ctx context.Context, role *models.UserRole, limit, offset int) ([]models.Profile, error) {
	var profiles []models.Profile
	query := `SELECT * FROM profiles WHERE ($1::user_role IS NULL OR role = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := database.DB.QueryContext(ctx, query, role, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var profile models.Profile
		if err := rows.Scan(&profile.ID, &profile.FullName, &profile.Phone, &profile.Role, &profile.Location, &profile.CreatedAt, &profile.UpdatedAt, &profile.Email); err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}
