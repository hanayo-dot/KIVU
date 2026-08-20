package models

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest payload for POST /auth/register.
type RegisterRequest struct {
	Name         string `json:"name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	LocationName string `json:"location_name"`
}

// LoginRequest payload for POST /auth/login.
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// RefreshTokenRequest payload for POST /auth/refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenResponse returned upon successful authentication.
type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	TokenType    string  `json:"token_type"` // Always "Bearer"
	ExpiresIn    int64   `json:"expires_in"` // Seconds
	Farmer       Farmer  `json:"farmer"`
}

// RefreshToken stored in DB (hashed).
type RefreshToken struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FarmerID  uuid.UUID `db:"farmer_id" json:"farmer_id"`
	TokenHash string    `db:"token_hash" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Revoked   bool      `db:"revoked" json:"revoked"`
}
