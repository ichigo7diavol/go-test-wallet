package models

import (
	"time"

	"github.com/google/uuid"
)

type WalletModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Balance   float32   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
