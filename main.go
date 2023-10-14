package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/dgraph-io/ristretto"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/storage"
	"github.com/parkerroan/go-api-badger-cache-tutorial/internal/webserver"
)

func main() {
	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure all paths cancel the context to release resources

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7, // number of keys to track frequency of (10M).
		MaxCost:     1e6, // maximum number of items (1M)
		BufferItems: 64,  // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	cacheClient := storage.NewCacheClient(cache)

	//New sqlx DB
	sqlxDB := sqlx.MustConnect("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	sqlClient := storage.NewSqlClient(sqlxDB)

	storageClient := storage.NewStorageClient(sqlClient, cacheClient)

	redis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if status := redis.Ping(ctx); status.Err() != nil {
		log.Fatalf("Redis Ping(): %v", status.Err())
	}

	streamName := "tutorial.inventory.items"
	streamProcessor := storage.NewStreamProcessor(streamName, redis)

	// Start stream processor in a goroutine
	go func() {
		if err := streamProcessor.Run(ctx, storageClient.ProcessCacheInvalidationMessage); err != nil {
			slog.Error(err.Error())
			// Cancel the context if there is an error and shut down the application.
			// 		Depending on your application and fault tolerance, you may want to
			// 		log and continue here instead of shutting down the application.
			cancel()
		}
	}()

	webserver := webserver.NewWebServer(storageClient)

	// Start HTTP server in a goroutine
	go func() {
		if err := webserver.Run(ctx); err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("ListenAndServe(): %v", err))
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
