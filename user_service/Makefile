include .env

.PHONY:build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/user_service ./cmd/user_service/main.go

.PHONY:run
run: build
	./bin/user_service

.PHONY:migrations-up
migrations-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations up

.PHONY:migrations-reset
migrations-reset:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations reset
