package repository

import (
	"github.com/sentrionic/valkyrie/model"
	"gorm.io/gorm"
)

// guildRepository is data/repository implementation
// of service layer UserRepository
type channelRepository struct {
	DB *gorm.DB
}

// NewChannelRepository is a factory for initializing User Repositories
func NewChannelRepository(db *gorm.DB) model.ChannelRepository {
	return &channelRepository{
		DB: db,
	}
}

func (r *channelRepository) Create(c *model.Channel) error {
	return r.DB.Create(&c).Error
}

func (r *channelRepository) GetGuildDefault(guildId string) (*model.Channel, error) {
	channel := model.Channel{}
	result := r.DB.
		Where("guild_id = ?", guildId).
		Order("created_at ASC").
		First(&channel)

	return &channel, result.Error
}
