package repository

import (
	"errors"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// guildRepository is data/repository implementation
// of service layer UserRepository
type guildRepository struct {
	DB *gorm.DB
}

// NewGuildRepository is a factory for initializing User Repositories
func NewGuildRepository(db *gorm.DB) model.GuildRepository {
	return &guildRepository{
		DB: db,
	}
}

func (r *guildRepository) List(uid string) (*[]model.GuildResponse, error) {
	var guilds []model.GuildResponse
	result := r.DB.Raw(`
		SELECT distinct g."id",
		g."name",
		g."owner_id",
		g."icon",
		g."created_at",
		g."updated_at",
		((SELECT c."last_activity"
		 FROM channels c
		 JOIN guilds g ON g.id = c."guild_id"
		 WHERE g.id = member."guild_id"
		 order by c."last_activity" DESC
		 limit 1) > member."last_seen") AS "hasNotification",
		(SELECT c.id AS "default_channel_id"
		FROM channels c
	    JOIN guilds g ON g.id = c."guild_id"
		WHERE g.id = member."guild_id"
		ORDER BY c."created_at"
		LIMIT 1)
		FROM guilds g
		JOIN members as member
		on g."id"::text = member."guild_id"
		WHERE member."user_id" = ?
		ORDER BY g."created_at";
	`, uid).Find(&guilds)

	return &guilds, result.Error
}

func (r *guildRepository) GuildMembers(userId string, guildId string) (*[]model.MemberResponse, error) {
	var members []model.MemberResponse
	result := r.DB.Raw(`
		SELECT u.id,
		u.username,
		u.image,
		u."is_online",
		u."created_at",
		u."updated_at",
		m.nickname,
		m.color,
		EXISTS(
			SELECT 1
			FROM users
			LEFT JOIN friends f ON users.id = f."user_id"
			WHERE f."friend_id" = u.id
			AND f."user_id" = ?
		) AS "isFriend"
		FROM users AS u
		JOIN members m ON u."id"::text = m."user_id"
		WHERE m."guild_id" = ?
		ORDER BY (CASE WHEN m.nickname notnull THEN m.nickname ELSE u.username END)
	`, userId, guildId).Find(&members)

	return &members, result.Error
}

func (r *guildRepository) Create(g *model.Guild) error {
	return r.DB.Debug().Create(&g).Error
}

func (r *guildRepository) FindUserByID(uid string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.
		Preload("Guilds").
		Where("id = ?", uid).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", uid)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *guildRepository) FindByID(id string) (*model.Guild, error) {
	guild := &model.Guild{}

	if err := r.DB.
		Preload(clause.Associations).
		Where("id = ?", id).
		First(&guild).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return guild, apperrors.NewNotFound("id", id)
		}
		return guild, apperrors.NewInternal()
	}

	return guild, nil
}

func (r *guildRepository) Save(g *model.Guild) error {
	return r.DB.Save(&g).Error
}
