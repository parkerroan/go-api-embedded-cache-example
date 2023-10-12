package storage

import (
	"context"
	"log"
	"log/slog"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/go-redis/redis/v8"
)

type StreamProcessor struct {
	StreamName string

	redis *redis.Client
	cache *badger.DB
}

func NewStreamProcessor(redis *redis.Client, cache *badger.DB) *StreamProcessor {
	return &StreamProcessor{
		redis: redis,
		cache: cache,
	}
}

// Run function to process messages from Redis one by one
func (p *StreamProcessor) Run(ctx context.Context, rdb *redis.Client) {
	// Start from the beginning of the stream
	id := "0"
	for {
		data, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{p.StreamName, id},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Println(err)
			return
		}

		for _, result := range data {
			for _, message := range result.Messages {
				if err := processMessage(message); err != nil {
					slog.Error(err.Error())

					//retry
					continue
				}
				// Update the last seen ID
				id = message.ID
			}
		}
		// Sleep for a short while before checking for new messages
		time.Sleep(1 * time.Second)
	}
}

// Process each message from the stream
func processMessage(message redis.XMessage) error {
	slog.Info(message.ID, message.Values)
	return nil
	// Add any additional processing logic here
}
