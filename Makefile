postgres:
	docker run --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d -p 5432:5432 postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root bank

dropdb:
	docker exec -it postgres dropdb bank

migrateup:
	migrate -path ./db/migration -path ./db/migration -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -path ./db/migration/ -database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc