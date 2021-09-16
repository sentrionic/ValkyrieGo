package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

func GetMockChannel(guildId string) *model.Channel {
	return &model.Channel{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		GuildID:      &guildId,
		Name:         RandStr(8),
		IsPublic:     true,
		IsDM:         false,
		LastActivity: time.Now(),
		PCMembers:    nil,
		Messages:     nil,
	}
}
