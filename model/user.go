package model

type User struct {
	BaseModel
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null;uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Image    string `json:"image"`
	IsOnline bool   `gorm:"default:true" json:"isOnline"`
}
