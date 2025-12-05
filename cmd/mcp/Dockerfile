FROM golang:alpine AS builder
WORKDIR /mcp
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /mcp/cmd/mcp
RUN go build -o mcp main.go

FROM alpine:latest AS runtime
WORKDIR /mcp
COPY --from=builder /mcp/cmd/mcp/mcp .
EXPOSE 6174
CMD ["./mcp"]
