FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/main cmd/server/main.go

FROM alpine:latest as runner

WORKDIR /app

COPY --from=builder /app/bin/main ./bin/main

EXPOSE 8080


CMD ["./bin/main"]