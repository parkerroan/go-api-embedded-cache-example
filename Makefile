LOCAL_PSQL_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

migrate.up:
	migrate -path migrations -database $(LOCAL_PSQL_URL) -verbose up

migrate.down:
	migrate -path migrations -database $(LOCAL_PSQL_URL) -verbose down 1

migrate.new:
	migrate create -ext sql -dir migrations $(name)

make db.reset:
	docker-compose down
	docker-compose up -d postgres
	sleep 10
	make migrate.up

make restart: 
	docker-compose down
	make db.reset
	docker-compose up -d redis debezium