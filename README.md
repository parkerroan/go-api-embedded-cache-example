# Go API Embedded Cache Example


# [Medium Article](https://medium.com/p/415ba52b5aa5)

This repo contains an API using Go that incorporates an embedded cache capable of scaling safely to any specified number of replicas.

In this pattern, we'll use the following components:

1. **Postgres** as the primary database.
2. **[Debezium](https://debezium.io/documentation/reference/stable/operations/debezium-server.html)** to capture changes in Postgres and push to Redis Streams.
3. **[Redis Streams](https://redis.io/docs/data-types/streams/#streams-basics)** to effectively handle and transmit change events.
4. **[Ristretto](https://github.com/dgraph-io/ristretto)** as an in-memory caching solution in Go.
5. **Go** as the primary language for implementing the service.


# Setup

## Dependencies

 * Golang 1.21.3
 * [Go Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

## Install Go Packages: 

 `go mod tidy`

## Run Locally

This will setup the db and services:

`make restart`

This will run the go application:

`go run main.go`

## Test It Out

[Postman Collection](https://www.postman.com/pkroan/workspace/blogs/collection/1491858-75528d19-b486-4281-bf01-0c6419c27c25?action=share&creator=1491858)
