package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	localCrypto "github.com/imrany/smart_spore_hub/server/database/crypto"
	"github.com/imrany/smart_spore_hub/server/database/models"
	"github.com/imrany/smart_spore_hub/server/database/processes/profile"
	"github.com/spf13/viper"
)

// CreateProfile creates a new user profile - POST /v1/profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")

	var userProfileRequest models.CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&userProfileRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	slog.Info("Creating new user profile", "data", userProfileRequest)

	userProfile := models.Profile{
		ID:        uuid.New().String(),
		FullName:  userProfileRequest.FullName,
		Email:     userProfileRequest.Email,
		Phone:     userProfileRequest.Phone,
		Role:      userProfileRequest.Role,
		Location:  userProfileRequest.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userProfileRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	userProfile.Password = string(hashedPassword)

	session, err := localCrypto.GenerateToken(userProfile.ID)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := profile.Create(ctx, &userProfile); err != nil {
		slog.Error("Failed to create profile", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	userProfile.Session = struct {
		ID        string `json:"id"`
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}{
		ID:        userProfile.ID,
		Token:     session,
		ExpiresAt: "1h",
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    userProfile,
	})
}

// LoginUser logs in a user - POST /v1/profile/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	userProfileRequest := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&userProfileRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	userProfile, err := profile.GetByEmail(ctx, userProfileRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userProfile.Password), []byte(userProfileRequest.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid email or password",
		})
		return
	}

	session, err := localCrypto.GenerateToken(userProfile.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	userProfile.Session = struct {
		ID        string `json:"id"`
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}{
		ID:        userProfile.ID,
		Token:     session,
		ExpiresAt: viper.GetString("JWT_EXPIRATION"),
	}
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    userProfile,
	})
}

// GetUserProfile retrieves the user profile information - GET /v1/profile/{id}
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	userProfile, err := profile.GetByID(ctx, userID)
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
		Data:    userProfile,
	})
}

// UpdateProfile updates the user profile information - PUT /v1/profile/{id}
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	userID := r.PathValue("id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	userProfileRequest := &models.UpdateProfileRequest{}
	if err := json.NewDecoder(r.Body).Decode(userProfileRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	userProfile, err := profile.GetByID(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userProfile.Password), []byte(userProfileRequest.OldPassword)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid password",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userProfileRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	userProfileUpdate := models.UpdateProfileRequest{
		FullName: userProfileRequest.FullName,
		Role:     userProfileRequest.Role,
		Location: userProfileRequest.Location,
		Phone:    userProfileRequest.Phone,
		Password: string(hashedPassword),
	}

	if err := profile.Update(ctx, userID, userProfileUpdate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    userProfile,
	})
}

// DeleteProfile deletes a user profile.
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Content-Type", "application/json")
	userID := r.PathValue("id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "User ID is required",
		})
		return
	}

	if err := profile.Delete(ctx, userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Profile deleted successfully",
	})
}
