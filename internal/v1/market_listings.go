package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/imrany/smart_spore_hub/server/database/models"
	marketlisting "github.com/imrany/smart_spore_hub/server/database/processes/market_listing"
)

// GetMarketListing retrieves a list of market listings. - GET /api/v1/market_listings
func GetMarketListing(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	filter := models.MarketListingFilter{
		Available: []bool{true}[0],
	}
	listing, err := marketlisting.List(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    listing,
	})
}

// GetMarketListingByID retrieves a market listing by ID. - GET /api/v1/market_listings/{id}
func GetMarketListingByID(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	listing, err := marketlisting.GetByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    listing,
	})
}

// CreateMarketListing creates a new market listing. - POST /api/v1/market_listings
func CreateMarketListing(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	var listing models.MarketListing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := marketlisting.Create(ctx, &listing); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    listing,
	})
}
