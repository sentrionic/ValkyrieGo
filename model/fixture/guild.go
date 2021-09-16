package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

func GetMockGuild(uid string) *model.Guild {
	ownerId := RandID()
	if uid != "" {
		ownerId = uid
	}

	return &model.Guild{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:    RandStr(8),
		OwnerId: ownerId,
	}
}

func GetGuildMember(guildId string) *model.Member {
	return &model.Member{
		UserID:    RandID(),
		GuildID:   guildId,
		Nickname:  nil,
		Color:     nil,
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
