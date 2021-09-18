package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

func GetMockChannel(guildId string) *model.Channel {

	var guild *string = nil
	if guildId != "" {
		guild = &guildId
	}

	return &model.Channel{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		GuildID:      guild,
		Name:         RandStr(8),
		IsPublic:     true,
		LastActivity: time.Now(),
	}
}
func GetMockDMChannel() *model.Channel {
	return &model.Channel{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:         RandID(),
		IsDM:         true,
		LastActivity: time.Now(),
	}
}
