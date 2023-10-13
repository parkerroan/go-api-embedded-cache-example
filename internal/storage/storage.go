package storage

import (
	"context"
	"log/slog"
	"time"
)

type Item struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type StorageClient struct {
	sql             *SqlClient
	cache           Cache
	streamProcessor *StreamProcessor
}

func NewStorageClient(sql *SqlClient, cache Cache) *StorageClient {

	storage := &StorageClient{
		sql:   sql,
		cache: cache,
	}

	return storage
}

func (s *StorageClient) GetItem(ctx context.Context, id string) (*Item, error) {
	cachedResult, err := s.cache.Retreive(id)
	if err != nil {
		return nil, err
	}

	if cachedResult != nil {
		return cachedResult, nil
	}

	result, err := s.sql.GetItem(ctx, id)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	if err := s.cache.Set(id, result); err != nil {
		//log error and continue
		slog.Error(err.Error())
	}

	return result, nil
}

func (s *StorageClient) UpsertItem(ctx context.Context, item Item) (*Item, error) {
	result, err := s.sql.UpsertItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *StorageClient) ProcessCacheInvalidationMessage(ctx context.Context, msgID string, values map[string]interface{}) error {
	//log message
	slog.Info("Processing cache invalidation message:", slog.String("msg_id", msgID))
	return s.cache.Remove(msgID)
}
