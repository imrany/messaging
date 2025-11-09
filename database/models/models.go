package models

import (
	"time"
)

// UserRole represents the enum for user roles
type UserRole string

const (
	UserRoleFarmer UserRole = "farmer"
	UserRoleBuyer  UserRole = "buyer"
	UserRoleAdmin  UserRole = "admin"
)

// Profile represents a user profile in the system
type Profile struct {
	ID        string    `json:"id" db:"id"`
	FullName  string    `json:"full_name" db:"full_name"`
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Email     string    `json:"email,omitempty" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
	Role      UserRole  `json:"role" db:"role"`
	Location  string    `json:"location,omitempty" db:"location"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Session   struct {
		ID        string `json:"id"`
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	} `json:"session"`
}

// Hub represents a physical hub location
type Hub struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Location     string    `json:"location" db:"location"`
	ManagerID    string    `json:"manager_id,omitempty" db:"manager_id"`
	Description  string    `json:"description,omitempty" db:"description"`
	ContactPhone string    `json:"contact_phone,omitempty" db:"contact_phone"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Course represents an educational course
type Course struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description,omitempty" db:"description"`
	Content     string    `json:"content,omitempty" db:"content"`
	ImageURL    string    `json:"image_url,omitempty" db:"image_url"`
	Duration    string    `json:"duration,omitempty" db:"duration"`
	Level       string    `json:"level" db:"level"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// MarketListing represents a product listing in the marketplace
type MarketListing struct {
	ID           string    `json:"id" db:"id"`
	FarmerID     string    `json:"farmer_id" db:"farmer_id"`
	ProductName  string    `json:"product_name" db:"product_name"`
	Description  string    `json:"description,omitempty" db:"description"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	Unit         string    `json:"unit" db:"unit"`
	PricePerUnit float64   `json:"price_per_unit" db:"price_per_unit"`
	Available    bool      `json:"available" db:"available"`
	ImageURL     string    `json:"image_url,omitempty" db:"image_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// NotificationPreference represents user notification settings
type NotificationPreference struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	SMSEnabled      bool      `json:"sms_enabled" db:"sms_enabled"`
	WhatsAppEnabled bool      `json:"whatsapp_enabled" db:"whatsapp_enabled"`
	EmailEnabled    bool      `json:"email_enabled" db:"email_enabled"`
	PhoneNumber     string    `json:"phone_number,omitempty" db:"phone_number"`
	Email           string    `json:"email,omitempty" db:"email"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// SensorReading represents a reading from IoT sensors
type SensorReading struct {
	ID          string    `json:"id" db:"id"`
	HubID       string    `json:"hub_id" db:"hub_id"`
	Temperature float64   `json:"temperature" db:"temperature"`
	Humidity    float64   `json:"humidity" db:"humidity"`
	RecordedAt  time.Time `json:"recorded_at" db:"recorded_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AlertType represents the type of alert
type AlertType string

const (
	AlertTypeTemperature AlertType = "temperature"
	AlertTypeHumidity    AlertType = "humidity"
	AlertTypeBoth        AlertType = "both"
)

// Alert represents a system alert
type Alert struct {
	ID          string    `json:"id" db:"id"`
	HubID       string    `json:"hub_id" db:"hub_id"`
	AlertType   AlertType `json:"alert_type" db:"alert_type"`
	Message     string    `json:"message" db:"message"`
	Temperature float64   `json:"temperature,omitempty" db:"temperature"`
	Humidity    float64   `json:"humidity,omitempty" db:"humidity"`
	Resolved    bool      `json:"resolved" db:"resolved"`
	ResolvedAt  time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ========================================
// Request/Response DTOs
// ========================================

// CreateProfileRequest represents the request to create a profile
type CreateProfileRequest struct {
	FullName string   `json:"full_name" validate:"required"`
	Phone    string   `json:"phone,omitempty"`
	Email    string   `json:"email,omitempty"`
	Role     UserRole `json:"role" validate:"required,oneof=farmer buyer admin"`
	Location string   `json:"location,omitempty"`
	Password string   `json:"password" validate:"required"`
}

// UpdateProfileRequest represents the request to update a profile
type UpdateProfileRequest struct {
	FullName    string   `json:"full_name,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Role        UserRole `json:"role,omitempty" validate:"omitempty,oneof=farmer buyer admin"`
	Location    string   `json:"location,omitempty"`
	Password    string   `json:"password,omitempty"`
	OldPassword string   `json:"old_password,omitempty"`
}

// CreateHubRequest represents the request to create a hub
type CreateHubRequest struct {
	Name         string `json:"name" validate:"required"`
	Location     string `json:"location" validate:"required"`
	ManagerID    string `json:"manager_id,omitempty"`
	Description  string `json:"description,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
}

// CreateCourseRequest represents the request to create a course
type CreateCourseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	Duration    string `json:"duration,omitempty"`
	Level       string `json:"level" validate:"required,oneof=beginner intermediate advanced"`
}

// CreateMarketListingRequest represents the request to create a market listing
type CreateMarketListingRequest struct {
	ProductName  string  `json:"product_name" validate:"required"`
	Description  string  `json:"description,omitempty"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	Unit         string  `json:"unit" validate:"required"`
	PricePerUnit float64 `json:"price_per_unit" validate:"required,gt=0"`
	ImageURL     string  `json:"image_url,omitempty"`
}

// UpdateMarketListingRequest represents the request to update a market listing
type UpdateMarketListingRequest struct {
	ProductName  string   `json:"product_name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Quantity     *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Unit         string   `json:"unit,omitempty"`
	PricePerUnit *float64 `json:"price_per_unit,omitempty" validate:"omitempty,gt=0"`
	Available    *bool    `json:"available,omitempty"`
	ImageURL     string   `json:"image_url,omitempty"`
}

// UpdateNotificationPreferenceRequest represents the request to update notification preferences
type UpdateNotificationPreferenceRequest struct {
	SMSEnabled      *bool  `json:"sms_enabled,omitempty"`
	WhatsAppEnabled *bool  `json:"whatsapp_enabled,omitempty"`
	EmailEnabled    *bool  `json:"email_enabled,omitempty"`
	PhoneNumber     string `json:"phone_number,omitempty"`
}

// CreateSensorReadingRequest represents the request to create a sensor reading
type CreateSensorReadingRequest struct {
	HubID       string    `json:"hub_id" validate:"required"`
	Temperature float64   `json:"temperature" validate:"required"`
	Humidity    float64   `json:"humidity" validate:"required"`
	RecordedAt  time.Time `json:"recorded_at" validate:"required"`
}

// CreateAlertRequest represents the request to create an alert
type CreateAlertRequest struct {
	HubID       string    `json:"hub_id" validate:"required"`
	AlertType   AlertType `json:"alert_type" validate:"required,oneof=temperature humidity both"`
	Message     string    `json:"message" validate:"required"`
	Temperature *float64  `json:"temperature,omitempty"`
	Humidity    *float64  `json:"humidity,omitempty"`
}

// ResolveAlertRequest represents the request to resolve an alert
type ResolveAlertRequest struct {
	Resolved bool `json:"resolved" validate:"required"`
}

// ========================================
// Query Filters
// ========================================

// SensorReadingFilter represents filters for querying sensor readings
type SensorReadingFilter struct {
	HubID     string    `json:"hub_id,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// AlertFilter represents filters for querying alerts
type AlertFilter struct {
	HubID     string    `json:"hub_id,omitempty"`
	AlertType AlertType `json:"alert_type,omitempty"`
	Resolved  bool      `json:"resolved,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// MarketListingFilter represents filters for querying market listings
type MarketListingFilter struct {
	FarmerID  string  `json:"farmer_id,omitempty"`
	Available bool    `json:"available,omitempty"`
	MinPrice  float64 `json:"min_price,omitempty"`
	MaxPrice  float64 `json:"max_price,omitempty"`
	Limit     int     `json:"limit,omitempty"`
	Offset    int     `json:"offset,omitempty"`
}
