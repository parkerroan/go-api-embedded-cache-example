# Go API Embedded Cache Example

This repo contains an API using Go that incorporates an embedded cache capable of scaling safely to any specified number of replicas.

The use of a embedded cache alongside change data capture (CDC) invalidation is a powerful technique offers a robust method to ensure our high-performance cache remains synchronized with its primary data source.

In this pattern, we'll use the following components:

1. **Postgres** as the primary database.
2. **[Debezium](https://debezium.io/documentation/reference/stable/operations/debezium-server.html)** to capture changes in Postgres and push to Redis Streams.
3. **[Redis Streams](https://redis.io/docs/data-types/streams/#streams-basics)** to effectively handle and transmit change events.
4. **[Ristretto](https://github.com/dgraph-io/ristretto)** as an in-memory caching solution in Go.
5. **Go** as the primary language for implementing the service.

# Diagrams: 

![Arch](https://substackcdn.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F199a014c-b55f-4211-923b-7c78f7a6735e_1502x1277.png)

![Read Sequence](https://substackcdn.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2F19adcfb7-59b4-4d1b-9cc0-477306ed834d_1200x647.png)

![Edit Sequence](https://substackcdn.com/image/fetch/f_auto,q_auto:good,fl_progressive:steep/https%3A%2F%2Fsubstack-post-media.s3.amazonaws.com%2Fpublic%2Fimages%2Fe7e69da0-3cd4-4175-bdb3-d853cd18afe1_800x725.png)


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
