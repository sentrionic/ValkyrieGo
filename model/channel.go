package model

import "time"

// Channel represents a text channel in a guild
// or a text channel for DMs between users.
// GuildID should only be nil if it is a DM channel
// PCMembers should only be used if the channel is private.
type Channel struct {
	BaseModel
	GuildID      *string
	Name         string
	IsPublic     bool
	IsDM         bool
	LastActivity time.Time `gorm:"autoCreateTime"`
	PCMembers    []User    `gorm:"many2many:pcmembers;constraint:OnDelete:CASCADE;"`
	Messages     []Message `gorm:"constraint:OnDelete:CASCADE;"`
}

// ChannelResponse is the JSON response of the channel
type ChannelResponse struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	IsPublic        bool      `json:"isPublic"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	HasNotification bool      `json:"hasNotification"`
} //@name Channel

// SerializeChannel returns the channel API response.
func (c Channel) SerializeChannel() ChannelResponse {
	return ChannelResponse{
		Id:              c.ID,
		Name:            c.Name,
		IsPublic:        c.IsPublic,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
		HasNotification: false,
	}
}
