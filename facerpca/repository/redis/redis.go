package redis

import (
	"context"
	"face-recognition-fyp/domain"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type FaceRPCARedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(client *redis.Client) domain.FaceRPCARedisRepository {
	return &FaceRPCARedisRepository{
		Client: client,
	}
}

func (r *FaceRPCARedisRepository) SetSession(ctx context.Context, token string, userID int, duration time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	return r.Client.Set(ctx, key, userID, duration).Err()
}

func (r *FaceRPCARedisRepository) GetSession(ctx context.Context, token string) (int, error) {
	key := fmt.Sprintf("session:%s", token)
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (r *FaceRPCARedisRepository) DeleteSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	return r.Client.Del(ctx, key).Err()
}
