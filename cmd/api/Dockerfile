FROM golang:alpine AS builder
WORKDIR /api
COPY go.mod go.sum ./
COPY . .
WORKDIR /api/cmd/api
RUN go build -o api main.go

FROM alpine:latest AS runtime
WORKDIR /api
COPY --from=builder /api/cmd/api/api .
EXPOSE 6173
CMD ["./api"]
