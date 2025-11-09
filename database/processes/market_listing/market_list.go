package marketlisting

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/imrany/smart_spore_hub/server/database"
	"github.com/imrany/smart_spore_hub/server/database/models"
)

func Create(ctx context.Context, listing *models.MarketListing) error {
	query := `
		INSERT INTO market_listings (farmer_id, product_name, description, quantity, unit, price_per_unit, available, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return database.DB.QueryRowContext(
		ctx, query,
		listing.FarmerID, listing.ProductName, listing.Description, listing.Quantity,
		listing.Unit, listing.PricePerUnit, listing.Available, listing.ImageURL,
	).Scan(&listing.ID, &listing.CreatedAt)
}

func GetByID(ctx context.Context, id string) (*models.MarketListing, error) {
	var listing models.MarketListing
	query := `SELECT id, farmer_id, product_name, description, quantity, unit, price_per_unit, available, image_url, created_at FROM market_listings WHERE id = $1`
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&listing.ID, &listing.FarmerID, &listing.ProductName, &listing.Description,
		&listing.Quantity, &listing.Unit, &listing.PricePerUnit, &listing.Available,
		&listing.ImageURL, &listing.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("market listing not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get market listing: %w", err)
	}
	return &listing, nil
}

func Update(ctx context.Context, id string, req models.UpdateMarketListingRequest) error {
	query := `
		UPDATE market_listings
		SET product_name = COALESCE($2, product_name),
		    description = COALESCE($3, description),
		    quantity = COALESCE($4, quantity),
		    unit = COALESCE($5, unit),
		    price_per_unit = COALESCE($6, price_per_unit),
		    available = COALESCE($7, available),
		    image_url = COALESCE($8, image_url)
		WHERE id = $1
	`
	_, err := database.DB.ExecContext(ctx, query, id, req.ProductName, req.Description, req.Quantity, req.Unit, req.PricePerUnit, req.Available, req.ImageURL)
	return err
}

func Delete(ctx context.Context, id string) error {
	query := `DELETE FROM market_listings WHERE id = $1`
	_, err := database.DB.ExecContext(ctx, query, id)
	return err
}

func List(ctx context.Context, filter models.MarketListingFilter) ([]models.MarketListing, error) {
	var listings []models.MarketListing
	query := `
		SELECT id, farmer_id, product_name, description, quantity, unit, price_per_unit, available, image_url, created_at FROM market_listings
		WHERE ($1::uuid IS NULL OR farmer_id = $1)
				AND ($2::boolean IS NULL OR available = $2)
				AND ($3::decimal IS NULL OR price_per_unit >= $3)
				AND ($4::decimal IS NULL OR price_per_unit <= $4)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := database.DB.QueryContext(ctx, query, filter.FarmerID, filter.Available, filter.MinPrice, filter.MaxPrice, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var listing models.MarketListing
		if err := rows.Scan(
			&listing.ID, &listing.FarmerID, &listing.ProductName, &listing.Description,
			&listing.Quantity, &listing.Unit, &listing.PricePerUnit, &listing.Available,
			&listing.ImageURL, &listing.CreatedAt,
		); err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listings, nil
}
