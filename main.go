package main

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/storage"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/webserver"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage := storage.NewStorageClient(db)

	webserver := webserver.NewWebServer(storage)
	webserver.Run()
}
