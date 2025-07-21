FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY common ../common

COPY vegetable_service/go.mod vegetable_service/go.sum .

RUN go mod download 

COPY vegetable_service/cmd cmd
COPY vegetable_service/internal internal
COPY vegetable_service/.env .env

RUN CGO_ENABLED=0 GOOS=linux go build -o vegetable_service ./cmd/vegetable_service/main.go

FROM alpine:3.21

COPY swagger swagger

# COPY vegetable_service/third_party third_party
COPY --from=builder /app/vegetable_service vegetable_service

# # Make sure the binary is executable
# RUN chmod +x ./vegetable_service

CMD [ "./vegetable_service" ]
# ENTRYPOINT ["./vegetable_service"]
