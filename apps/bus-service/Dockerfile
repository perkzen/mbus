FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/main cmd/server/main.go

RUN go build -o ./bin/seed cmd/seed/main.go

# ───────────────────────────────────────────────
FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/bin/seed ./bin/seed

EXPOSE 8080

CMD ["sh", "-c", "./bin/seed && ./bin/main"]
