FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY common ../common
COPY user_service/go.mod user_service/go.sum .

RUN go mod download 

COPY user_service/cmd cmd
COPY user_service/internal internal
COPY user_service/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -o user_service ./cmd/user_service/main.go

FROM alpine:3.21

COPY swagger swagger

COPY --from=builder /app/user_service user_service

CMD [ "./user_service" ]
