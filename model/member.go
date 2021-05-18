package model

import "time"

type Member struct {
	BaseModel
	UserID   string    `gorm:"primaryKey"`
	GuildID  string    `gorm:"primaryKey;"`
	Nickname *string   `gorm:"nickname"`
	Color    *string   `gorm:"color"`
	LastSeen time.Time `gorm:"autoCreateTime"`
}

type MemberResponse struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Image     string    `json:"image"`
	IsOnline  bool      `json:"isOnline"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Nickname  *string   `json:"nickname"`
	Color     *string   `json:"color"`
	IsFriend  bool      `json:"isFriend"`
}
