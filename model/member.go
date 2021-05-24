package model

import "time"

type Member struct {
	UserID    string    `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	GuildID   string    `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	Nickname  *string   `gorm:"nickname"`
	Color     *string   `gorm:"color"`
	LastSeen  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
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

type BanResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
}

type MemberSettings struct {
	Nickname *string `json:"nickname" binding:"omitempty,gte=3,lte=30"`
	Color    *string `json:"color" binding:"omitempty,hexcolor"`
}
