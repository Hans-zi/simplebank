startdb:
	docker start postgres14

postgres:
	docker run --name postgres14 -e POSTGRES_USER=root -e POSTGRES_PASSWORD -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	 migrate -path db/migration -database "postgresql://root:xiaohan1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgresql://root:xiaohan1234@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v $(CURDIR):/src -w /src sqlc/sqlc generate

test:
	go test -v -cover -short ./...

.PHONY: startdb postgres createdb dropdb migrateup migratedown sqlc test