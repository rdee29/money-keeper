package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid"`
	Amount      float64
	Type        string // income / expense
	Description string
	CreatedAt   time.Time
}