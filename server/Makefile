DB_URL=postgresql://root:secret@localhost:5432/carrot_market?sslmode=disable

network:
	docker network create carrot-market-network

postgres:
	docker run --name postgres -p 5432:5432 --network carrot-market-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root carrot_market

dropdb:
	docker exec -it postgres dropdb carrot_market

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

server:
	go run main.go

sqlc:
	sqlc generate
	mockgen -package mockdb -destination db/mock/store.go github.com/gitaepark/carrot-market/db/sqlc Store

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/gitaepark/carrot-market/db/sqlc Store

test:
	go test -v -cover ./...

.PHONY: network postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 db_docs db_schema server sqlc mock test