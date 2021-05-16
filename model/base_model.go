package model

import "time"

type BaseModel struct {
	ID        string    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"default:now();index" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updatedAt"`
}
