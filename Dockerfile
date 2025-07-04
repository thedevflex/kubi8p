FROM golang:1.24.1 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kubi8p ./cmd/main.go

# Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/kubi8p /usr/local/bin/kubi8p
COPY --from=builder /app/public /app/public

ENV PUBLIC_DIR=/app/public

ENTRYPOINT ["/usr/local/bin/kubi8p"]