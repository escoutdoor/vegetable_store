FROM golang:1.23.3-alpine3.21 as builder
WORKDIR /app

COPY ["go.mod", "go.sum"] ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/order_service ./cmd/order_service/main.go

FROM alpine:3.21
COPY --from=builder /app/bin/order_service /app

CMD [ "/app/order_service" ]
