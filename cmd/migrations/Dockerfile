FROM golang:alpine AS builder
WORKDIR /migrations
COPY go.mod go.sum ./
COPY . .
WORKDIR /migrations/cmd/migrations
RUN go build -o migrations_runner main.go

FROM alpine:latest AS runtime
WORKDIR /migrations
COPY --from=builder /migrations/cmd/migrations/migrations_runner .
COPY --from=builder /migrations/cmd/migrations/migrations ./migrations
CMD ["./migrations_runner"]
