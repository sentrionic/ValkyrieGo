package model

// User represents the user of the website.
type User struct {
	BaseModel
	Username string    `gorm:"not null" json:"username"`
	Email    string    `gorm:"not null;uniqueIndex" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	Image    string    `json:"image"`
	IsOnline bool      `gorm:"default:true" json:"isOnline"`
	Friends  []User    `gorm:"many2many:friends;" json:"-"`
	Requests []User    `gorm:"many2many:friend_requests;joinForeignKey:sender_id;joinReferences:receiver_id" json:"-"`
	Guilds   []Guild   `gorm:"many2many:members;" json:"-"`
	Message  []Message `json:"-"`
}
