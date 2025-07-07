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

run: 
	go run main.go

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown migratedown1 migrateforce sqlc docker-compose run