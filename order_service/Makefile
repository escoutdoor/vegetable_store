EXPORT POSTGRES_DSN=postgres://escoutdoor:escoutdoor@localhost:5412/orders

.PHONY:build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/order_service ./cmd/order_service/main.go

.PHONY:run
run: build
	./bin/order_service

.PHONY:migrations-up
migrations-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations up

.PHONY:migrations-down
migrations-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations down

.PHONY:migrations-reset
migrations-reset:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=${POSTGRES_DSN} goose -dir=migrations reset
