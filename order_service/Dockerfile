FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY common ../common

COPY order_service/go.mod order_service/go.sum .

RUN go mod download 

COPY order_service/cmd cmd
COPY order_service/internal internal
COPY order_service/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -o order_service ./cmd/order_service/main.go

FROM alpine:3.21

COPY swagger swagger

COPY --from=builder /app/order_service order_service

CMD [ "./order_service" ]
