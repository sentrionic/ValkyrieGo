package model

import "time"

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

type ChannelResponse struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	IsPublic        bool      `json:"isPublic"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	HasNotification bool      `json:"hasNotification"`
}

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
