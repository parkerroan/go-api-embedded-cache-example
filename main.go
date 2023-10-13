package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-redis/redis/v8"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/storage"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/webserver"
)

func main() {
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure all paths cancel the context to release resources

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cacheClient := storage.NewCacheClient(db)

	storageClient := storage.NewStorageClient(cacheClient)

	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	streamName := "tutorial.inventory.items"
	streamProcessor := storage.NewStreamProcessor(streamName, redis)

	// Start stream processor in a goroutine
	go func() {
		if err := streamProcessor.Run(ctx, storageClient.ProcessCacheInvalidationMessage); err != nil {
			slog.Error(err.Error())
			cancel()
		}
	}()

	webserver := webserver.NewWebServer(storageClient)

	// Start HTTP server in a goroutine
	go func() {
		if err := webserver.Run(ctx); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
			cancel()
		}
	}()

	// Listen for the cancellation
	select {
	case <-ctx.Done():
		slog.Info("Context was canceled. Exiting...")
		return
	}

}
