package model

import "time"

// Message represents a text message in a channel.
// It may contain an Attachment that is displayed instead of text.
type Message struct {
	BaseModel
	Text       *string
	UserId     string      `gorm:"constraint:OnDelete:CASCADE;"`
	ChannelId  string      `gorm:"constraint:OnDelete:CASCADE;"`
	Attachment *Attachment `gorm:"constraint:OnDelete:CASCADE;"`
}

// MessageResponse is the API response of a Message
type MessageResponse struct {
	Id         string         `json:"id"`
	Text       *string        `json:"text"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Attachment *Attachment    `json:"attachment"`
	User       MemberResponse `json:"user"`
}

// Attachment represents a message attachment that displays
// a file instead of text.
type Attachment struct {
	ID        string    `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time `gorm:"index" json:"-"`
	UpdatedAt time.Time `json:"-"`
	Url       string    `json:"url"`
	FileType  string    `json:"filetype"`
	Filename  string    `json:"filename"`
	MessageId string    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
