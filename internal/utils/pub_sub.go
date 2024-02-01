package utils

import (
	"concurrent-image-fetcher/internal/responses/ws"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v9"
)

// PubSub represents a Redis client with publish and subscribe capabilities.
type (
	PubSubInterface interface {
		PublishProgress(ctx context.Context, channel string, progress ws.DownloadProgress) error
	}
	PubSub struct {
		redis *redis.Client
	}
)

// PublishProgress publishes the download progress to a specific channel in Redis.
func (c *PubSub) PublishProgress(ctx context.Context, channel string, progress ws.DownloadProgress) error {
	message, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("error encoding progress: %w", err)
	}

	if err := c.redis.Publish(ctx, channel, message).Err(); err != nil {
		return fmt.Errorf("error publishing progress: %w", err)
	}

	return nil
}

// NewPubSub creates and returns a new instance of a Redis client.
func NewPubSub(redis *redis.Client) PubSubInterface {
	return &PubSub{
		redis: redis,
	}
}
