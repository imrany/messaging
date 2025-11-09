package course

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, course *models.Course) error {
	query := `
		INSERT INTO courses (title, description, content, image_url, duration, level)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		course.Title, course.Description, course.Content, course.ImageURL, course.Duration, course.Level,
	).Scan(&course.ID, &course.CreatedAt)
}

func GetByID(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	query := `SELECT id, title, description, content, image_url, duration, level, created_at FROM courses WHERE id = $1`
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&course.ID, &course.Title, &course.Description, &course.Content,
		&course.ImageURL, &course.Duration, &course.Level, &course.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("course not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}
	return &course, nil
}

func Update(ctx context.Context, id string, course *models.Course) error {
	query := `
		UPDATE courses
		SET title = $2, description = $3, content = $4, image_url = $5, duration = $6, level = $7
		WHERE id = $1
	`
	_, err := database.DB.ExecContext(ctx, query, id, course.Title, course.Description, course.Content, course.ImageURL, course.Duration, course.Level)
	return err
}

func Delete(ctx context.Context, id string) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := database.DB.ExecContext(ctx, query, id)
	return err
}

func List(ctx context.Context, level *string, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	query := `SELECT id, title, description, content, image_url, duration, level, created_at FROM courses WHERE ($1::text IS NULL OR level = $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := database.DB.QueryContext(ctx, query, level, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course models.Course
		if err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.Content,
			&course.ImageURL, &course.Duration, &course.Level, &course.CreatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
