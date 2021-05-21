package model

import "time"

type DMMember struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"primaryKey"`
	ChannelId string `gorm:"primaryKey;"`
	IsOpen    bool
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}
