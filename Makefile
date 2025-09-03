DB_USER=root
DB_PASSWORD=secret
DB_CONTAINER=postgres13
DB_NAME=simple_bank

postgres:
	docker run --name $(DB_CONTAINER) -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:13-alpine

createdb:
	docker exec -it $(DB_CONTAINER) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER) dropdb $(DB_NAME)

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test