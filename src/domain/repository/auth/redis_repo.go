package auth_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/andy74139/webserver/src/database/redis"
	"github.com/andy74139/webserver/src/domain/entity/auth"
)

// postgresql account repository
type redisRepo struct {
	cache *redis.Client
}

func NewRedisRepo(rdb *redis.Client) (auth.Repository, error) {
	if rdb == nil {
		return nil, errors.New("redis repository can't be nil")
	}
	return &redisRepo{cache: rdb}, nil
}

func (r *redisRepo) SetRevoked(ctx context.Context, jwtID string, expiryDuration time.Duration) error {
	// TODO: cache-aside with db, to keep revoking token if cache being down
	key := getAuthTokenRedisKey(jwtID)
	if err := r.cache.Set(ctx, key, nil, expiryDuration).Err(); err != nil {
		return fmt.Errorf("Set error: %w", err)
	}
	return nil
}

func (r *redisRepo) IsRevoked(ctx context.Context, jwtID string) (bool, error) {
	if err := r.cache.Get(ctx, getAuthTokenRedisKey(jwtID)).Err(); errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("IsRevoked error: %w", err)
	}
	return true, nil
}

func getAuthTokenRedisKey(tokenID string) string {
	return rediskey.GetRevokeAuthPrefix() + tokenID
}
