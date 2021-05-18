package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sentrionic/valkyrie/model"
	_ "image/jpeg"
	_ "image/png"
	"time"
)

type redisRepository struct {
	rds *redis.Client
}

// NewRedisRepository is a factory for initializing User Repositories
func NewRedisRepository(rds *redis.Client) model.RedisRepository {
	return &redisRepository{
		rds: rds,
	}
}

var (
	InviteLinkPrefix     = "inviteLink"
	ForgotPasswordPrefix = "forgot-password"
)

func (r *redisRepository) SetResetToken(ctx context.Context, id string) (string, error) {
	uid, err := gonanoid.New()

	if err != nil {
		return "", err
	}

	if err := r.rds.Set(ctx, fmt.Sprintf("%s:%s", ForgotPasswordPrefix, uid), id, 24*time.Hour).Err(); err != nil {
		fmt.Println(err)
		return "", err
	}

	return uid, nil
}

func (r *redisRepository) GetIdFromToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", ForgotPasswordPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	r.rds.Del(ctx, key)

	return val, nil
}

func (r *redisRepository) SaveInvite(ctx context.Context, guildId string, id string, isPermanent bool) error {

	invite := model.Invite{GuildId: guildId, IsPermanent: isPermanent}

	value, err := json.Marshal(invite)

	if err != nil {
		return err
	}

	expiration := 24 * time.Hour
	if isPermanent {
		expiration = 0
	}

	result := r.rds.Set(ctx, fmt.Sprintf("%s:%s", InviteLinkPrefix, id), value, expiration)
	return result.Err()
}

func (r *redisRepository) GetInvite(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", InviteLinkPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var invite model.Invite
	err = json.Unmarshal([]byte(val), &invite)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if !invite.IsPermanent {
		r.rds.Del(ctx, key)
	}

	return invite.GuildId, nil
}

func (r *redisRepository) InvalidateInvites(ctx context.Context, guild *model.Guild) {
	for _, v := range guild.InviteLinks {
		key := fmt.Sprintf("%s:%s", InviteLinkPrefix, v)
		r.rds.Del(ctx, key)
	}
}
