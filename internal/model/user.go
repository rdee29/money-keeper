package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name 		string
	Email 		string `gorm:"unique"`
	Password 	string
	CreatedAt 	time.Time
}