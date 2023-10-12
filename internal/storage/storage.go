package storage

import (
	"encoding/json"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

type Item struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

type StorageClient struct {
	Cache *badger.DB
}

func NewStorageClient(cache *badger.DB) *StorageClient {
	return &StorageClient{
		Cache: cache,
	}
}

func (s *StorageClient) GetItem(id string) (*Item, error) {
	cachedResult, err := s.RetreiveFromCache(id)
	if err != nil {
		return nil, err
	}

	if cachedResult != nil {
		return cachedResult, nil
	}

	//TODO Call DB to get item

	return &Item{}, nil
}

func (s *StorageClient) UpsertItem(id string) (*Item, error) {
	return &Item{}, nil
}

func (s *StorageClient) RetreiveFromCache(id string) (*Item, error) {
	txn := s.Cache.NewTransaction(true)
	defer txn.Discard()

	item, err := txn.Get([]byte(id))
	if err != nil {
		return nil, err
	}

	// take item bytes and unmarshal into Item struct
	var value []byte
	value, err = item.ValueCopy(value)
	if err != nil {
		return nil, err
	}

	//marshal item bytes into Item struct
	var result Item
	err = json.Unmarshal(value, &result)

	return &result, nil
}
