package storage

import (
	"encoding/json"

	badger "github.com/dgraph-io/badger/v4"
)

type Cache interface {
	Set(id string, item *Item) error
	Retreive(id string) (*Item, error)
	Remove(id string) error
}

type CacheClient struct {
	db *badger.DB
}

func NewCacheClient(db *badger.DB) *CacheClient {
	return &CacheClient{
		db: db,
	}
}

func (c *CacheClient) Set(id string, item *Item) error {
	return nil
}

func (c *CacheClient) Retreive(id string) (*Item, error) {
	txn := c.db.NewTransaction(true)
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

	if err := txn.Commit(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *CacheClient) Remove(id string) error {
	txn := c.db.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Delete([]byte(id)); err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}
