package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"strconv"
	"time"
)

var (
	ErrItemNotFound = errors.New("item not found")
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
	sql   *SqlClient
	cache Cache
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
		slog.Info("Retrieved item from cache", slog.String("id", cachedResult.ID))
		return cachedResult, nil
	}

	result, err := s.sql.GetItem(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	//set item in cache
	if err := s.cache.Set(id, *result); err != nil {
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

	type msg struct {
		Payload struct {
			Op    string `json:"op"`
			After struct {
				ID int `json:"id"`
			} `json:"after"`
		} `json:"payload"`
	}

	//extract payload
	for _, v := range values {
		//marshal to Payload struct
		v := v.(string)
		var message msg
		if err := json.Unmarshal([]byte(v), &message); err != nil {
			return err
		}
		if message.Payload.Op != "" {
			//remove item from cache
			slog.Info("Removing item from cache:", slog.Int("id", message.Payload.After.ID))

			//convert ID to string
			id := strconv.Itoa(message.Payload.After.ID)
			if err := s.cache.Remove(id); err != nil {
				return err
			}
		}
	}

	return nil
}
