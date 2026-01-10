package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"user-restriction-manager/internal/core/domain"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) CacheActiveRestrictions(ctx context.Context, userID string, restrictions []domain.Restriction) error {
	data, err := json.Marshal(restrictions)
	if err != nil {
		return fmt.Errorf("failed to marshal restrictions: %w", err)
	}

	key := fmt.Sprintf("restrictions:%s", userID)
	// Cache for 15 minutes, or until invalidated
	err = r.client.Set(ctx, key, data, 15*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to cache restrictions: %w", err)
	}
	return nil
}

func (r *CacheRepository) GetActiveRestrictions(ctx context.Context, userID string) ([]domain.Restriction, error) {
	key := fmt.Sprintf("restrictions:%s", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cached restrictions: %w", err)
	}

	var restrictions []domain.Restriction
	if err := json.Unmarshal([]byte(val), &restrictions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal restrictions: %w", err)
	}

	return restrictions, nil
}

func (r *CacheRepository) Invalidate(ctx context.Context, userID string) error {
	key := fmt.Sprintf("restrictions:%s", userID)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}
	return nil
}
