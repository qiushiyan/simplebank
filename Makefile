.PHONY: createdb dropdb

pg:
	docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5433:5432 -d bank-api-postgres

pgcli:
	docker exec -it bank-api-postgres psql -U postgres

createdb:
	docker exec -it bank-api-postgres createdb --username=postgres --owner=postgres bank

dropdb:
	docker exec -it bank-api-postgres dropdb bank --username=postgres

migrate_create:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrate_up:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable" --verbose up

migrate_down:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable" --verbose down

generate:
	sqlc generate

test:
	go test -v -cover ./...

check:
	nilaway app/*