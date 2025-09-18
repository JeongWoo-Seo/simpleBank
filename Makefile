DB_USER=root
DB_PASSWORD=secret
DB_CONTAINER=postgres13
DB_NAME=simple_bank
DB_URL=postgres://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name $(DB_CONTAINER) --network simplebank-network -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:13-alpine

createdb:
	docker exec -it $(DB_CONTAINER) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER) dropdb $(DB_NAME)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb \
  -destination db/mock/store.go \
  github.com/JeongWoo-Seo/simpleBank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock dbdocs dbml2sql