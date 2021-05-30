package model

import (
	"time"
)

// BaseModel is similar to the gorm.Model and it includes the
// ID as a string, CreatedAt and UpdatedAt fields
type BaseModel struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
