package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/service"
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

func (r *redisRepository) SetResetToken(ctx context.Context, id string) (string, error) {
	uid, err := service.GenerateId()

	if err != nil {
		return "", err
	}

	if err := r.rds.Set(ctx, fmt.Sprintf("forgot-password:%s", uid), id, 24*time.Hour).Err(); err != nil {
		fmt.Println(err)
		return "", err
	}

	return uid, nil
}

func (r *redisRepository) GetIdFromToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("forgot-password:%s", token)
	val, err := r.rds.Get(ctx, key).Result()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	r.rds.Del(ctx, key)

	return val, nil
}
