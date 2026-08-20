package models

import (
	"time"

	"github.com/google/uuid"
)

// Farmer represents an aquaculture cage farmer entity.
type Farmer struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	PhoneNumber  string    `db:"phone_number" json:"phone_number"`
	PasswordHash string    `db:"password_hash" json:"-"`
	LocationName string    `db:"location_name" json:"location_name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
