include app.env
export $(shell sed 's/=.*//' app.env)

postgres:
	docker run --name postgres -p $(POSTGRES_PORT):5432 \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-d postgres:15
createdb:
	docker exec -it postgres createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

dropdb:
	docker exec -it postgres dropdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

migrateup:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" -verbose down 1

migrateforce:
	migrate -path ./internal/db/migrations -database "$(DB_URL)" -verbose force 1

sqlc:
	sqlc generate

new-migration:
	migrate create -ext sql -dir ./internal/db/migrations -seq drop_admin

docker-compose:
	docker-compose --env-file app.env up -d

PROTO_FILES=$(shell find proto -name "*.proto" | grep -v "proto/google" | grep -v "proto/protoc-gen-openapiv2")

proto:
	rm -rf pb/*
	rm -rf doc/swagger/*.swagger.json
	rm -rf doc/statik/*
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=join_love \
	$(PROTO_FILES)
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

run: 
	go run cmd/main.go

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown migratedown1 migrateforce sqlc docker-compose proto evans run