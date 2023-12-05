DB_URL=postgresql://root:secret@localhost:5432/mediahls?sslmode=disable

network:
	docker network create media-hls-network

postgres:
	docker run --name postgres --network media-hls-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root mediahls

dropdb:
	docker exec -it postgres dropdb mediahls

psql:
	docker exec -it postgres psql -U root -d mediahls

migrateinit:
	migrate create -ext sql -dir internal/db/schema -seq init_schema

migrateup:
	migrate -path internal/db/schema -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path internal/db/schema -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: network postgres createdb dropdb psql migrateup migratedown sqlc test server