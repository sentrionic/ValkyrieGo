package model

import "time"

type Channel struct {
	BaseModel
	GuildID      string
	Name         string    `gorm:"not null"`
	IsPublic     bool      `gorm:"default:true"`
	IsDM         bool      `gorm:"default:false"`
	LastActivity time.Time `gorm:"autoCreateTime"`
	Members      []User    `gorm:"many2many:channel_members;joinForeignKey:channels;joinReferences:users"`
}
