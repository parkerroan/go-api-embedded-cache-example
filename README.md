# Go API Embedded Cache Example

This repo contains an API using Go that incorporates an embedded cache capable of scaling safely to any specified number of replicas.

The use of a embedded cache alongside change data capture (CDC) invalidation is a powerful technique offers a robust method to ensure our high-performance cache remains synchronized with its primary data source.

In this pattern, we'll use the following components:

1. **Postgres** as the primary database.
2. **[Debezium](https://debezium.io/documentation/reference/stable/operations/debezium-server.html)** to capture changes in Postgres and push to Redis Streams.
3. **[Redis Streams](https://redis.io/docs/data-types/streams/#streams-basics)** to effectively handle and transmit change events.
4. **[Ristretto](https://github.com/dgraph-io/ristretto)** as an in-memory caching solution in Go.
5. **Go** as the primary language for implementing the service.

# Sequence Diagrams: 

![Read Sequence](https://cdn-images-1.medium.com/max/800/1*wmgCh0COc--30vBW6Xdi_g.png)

![Edit Sequence](https://cdn-images-1.medium.com/max/1200/1*drE4ENgJ9WbnFhs4ulQbBA.png)


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
