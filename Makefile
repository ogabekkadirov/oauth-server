include .env

SERVICE={service-name}
MIGRATION_PATH=./src/infrastructure/db/migrations
PROTOS_PATH=./src/infrastructure/protos
PROTO_DIRS={proto file dir name}

server:
	go run main.go

seed:
	go run ./src/Infrastructure/cmd/seeder/main.go

migration:
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq ... ... ... $(table)

migrateup:
	migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable&search_path=public" \
	-path ${MIGRATION_PATH} up

migratedown:
	migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable&search_path=public" \
	-path ${MIGRATION_PATH} down


clear:
	rm -rf ./src/application/protos/*

compose-up:
	docker-compose -f ./deploy/docker-compose.yml up

compose-down:
	docker-compose -f ./deploy/docker-compose.yml down

docker:
	docker build --rm -t oauth-svc -f ./deploy/Dockerfile .