package redis

import (
	"context"
	"ecommerce-backend/internal/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartCacheRepository struct {
	client *redis.Client
}

func NewCartCacheRepository(client *redis.Client) *CartCacheRepository {
	return &CartCacheRepository{client: client}
}

func (c *CartCacheRepository) GetCart(ctx context.Context, userID int64) (*domain.CartSummaryResponse, error) {
	if c.client == nil {
		return nil, nil
	}

	key := fmt.Sprintf("cache:cart:user:%d", userID)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var summary domain.CartSummaryResponse
	if err := json.Unmarshal([]byte(val), &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func (c *CartCacheRepository) SetCart(ctx context.Context, userID int64, summary *domain.CartSummaryResponse, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	key := fmt.Sprintf("cache:cart:user:%d", userID)
	bytes, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, bytes, ttl).Err()
}

func (c *CartCacheRepository) InvalidateCart(ctx context.Context, userID int64) error {
	if c.client == nil {
		return nil
	}

	key := fmt.Sprintf("cache:cart:user:%d", userID)
	return c.client.Del(ctx, key).Err()
}
