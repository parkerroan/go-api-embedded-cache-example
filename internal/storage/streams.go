package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
)

// StreamProcessor is a struct that processes messages from Redis
type StreamProcessor struct {
	StreamName string

	rdb   *redis.Client
	cache Cache
}

// ProccessMessageFunc is a function that processes a message from Redis
type ProccessMessageFunc func(ctx context.Context, id string, values map[string]interface{}) error

func NewStreamProcessor(streamName string, redis *redis.Client) *StreamProcessor {
	return &StreamProcessor{
		StreamName: streamName,
		rdb:        redis,
	}
}

// Run function to process messages from Redis one by one
func (p *StreamProcessor) Run(ctx context.Context, processMessage ProccessMessageFunc) error {
	slog.Info("Starting stream processor for stream:", slog.String("stream_name", p.StreamName))
	// Start from the end of the stream
	id := "$"
	for {
		data, err := p.rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{p.StreamName, id},
			Count:   4,
			Block:   100,
		}).Result()
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		for _, result := range data {
			for _, message := range result.Messages {
				msg := message
				if err := processMessage(ctx, msg.ID, msg.Values); err != nil {
					slog.Error(err.Error())

					//retry
					continue
				}
				// Update the last seen ID
				id = msg.ID
			}
		}
		// Sleep for a short while before checking for new messages
		time.Sleep(50 * time.Millisecond)
	}
}
