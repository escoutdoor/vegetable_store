# include .env

export POSTGRES_DSN=postgres://escoutdoor:escoutdoor@localhost:5432/vegetables

.PHONY:build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/vegetable_service ./cmd/vegetable_service/main.go

.PHONY:run
run: build
	./bin/vegetable_service

.PHONY:migrations-up
migrations-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations up

.PHONY:migrations-down
migrations-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations down
