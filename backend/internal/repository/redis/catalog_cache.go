package redis

import (
	"context"
	"crypto/md5"
	"ecommerce-backend/internal/domain"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CatalogCacheRepository struct {
	client *redis.Client
}

func NewCatalogCacheRepository(client *redis.Client) *CatalogCacheRepository {
	return &CatalogCacheRepository{client: client}
}

type CachedProductList struct {
	Products []domain.Product `json:"products"`
	Total    int64            `json:"total"`
}

func (c *CatalogCacheRepository) GenerateProductFilterKey(filter domain.ProductFilter) string {
	raw := fmt.Sprintf("p:%d:l:%d:s:%s:st:%s", filter.Page, filter.Limit, filter.Search, filter.Status)
	if filter.CategoryID != nil {
		raw += fmt.Sprintf(":c:%d", *filter.CategoryID)
	}
	if filter.Sort != "" {
		raw += fmt.Sprintf(":srt:%s", filter.Sort)
	}
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("cache:products:list:%s", hex.EncodeToString(hash[:8]))
}

func (c *CatalogCacheRepository) GetProducts(ctx context.Context, key string) (*CachedProductList, error) {
	if c.client == nil {
		return nil, nil
	}

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cached CachedProductList
	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		return nil, err
	}

	return &cached, nil
}

func (c *CatalogCacheRepository) SetProducts(ctx context.Context, key string, products []domain.Product, total int64, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	cached := CachedProductList{
		Products: products,
		Total:    total,
	}

	bytes, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, bytes, ttl).Err()
}

func (c *CatalogCacheRepository) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	if c.client == nil {
		return nil, nil
	}

	key := fmt.Sprintf("cache:products:slug:%s", slug)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var product domain.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *CatalogCacheRepository) SetProductBySlug(ctx context.Context, slug string, product *domain.Product, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	key := fmt.Sprintf("cache:products:slug:%s", slug)
	bytes, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, bytes, ttl).Err()
}

func (c *CatalogCacheRepository) InvalidateProductCache(ctx context.Context) error {
	if c.client == nil {
		return nil
	}

	// Delete all keys matching cache:products:*
	iter := c.client.Scan(ctx, 0, "cache:products:*", 0).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
	}
	return iter.Err()
}
