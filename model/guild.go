package model

import (
	"github.com/lib/pq"
	"time"
)

type Guild struct {
	BaseModel
	Name        string `gorm:"not null"`
	OwnerId     string `gorm:"not null"`
	Icon        *string
	InviteLinks pq.StringArray `gorm:"type:text[]"`
	Members     []User         `gorm:"many2many:members;constraint:OnDelete:CASCADE;"`
	Channels    []Channel      `gorm:"constraint:OnDelete:CASCADE;"`
	Bans        []User         `gorm:"many2many:bans;constraint:OnDelete:CASCADE;"`
}

type GuildResponse struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	OwnerId          string    `json:"ownerId"`
	Icon             *string   `json:"icon"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	HasNotification  bool      `json:"hasNotification"`
	DefaultChannelId string    `json:"default_channel_id"`
}

// SerializeGuild returns the guild API response.
// The channelId represents the default channel the user gets send to.
func (g Guild) SerializeGuild(channelId string) GuildResponse {
	return GuildResponse{
		Id:               g.ID,
		Name:             g.Name,
		OwnerId:          g.OwnerId,
		Icon:             g.Icon,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		HasNotification:  false,
		DefaultChannelId: channelId,
	}
}
