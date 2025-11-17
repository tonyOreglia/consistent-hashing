package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client
type Client struct {
	client *redis.Client
	ctx    context.Context
}

// Config holds Redis connection configuration
// New creates a new Redis client from a URL (assumes no password)
func New(url string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // No password assumed
		DB:       0,  // just connect to DB 0
	})

	ctx := context.Background()

	// Test the connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Set stores a string value by key
func (c *Client) Set(key, value string) error {
	return c.client.Set(c.ctx, key, value, 0).Err()
}

// Get retrieves a string value by key
func (c *Client) Get(key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist: %s", key)
	} else if err != nil {
		return "", err
	}
	return val, nil
}

// Delete removes a key from Redis
func (c *Client) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists checks if a key exists
func (c *Client) Exists(key string) (bool, error) {
	n, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.client.Close()
}
