# syntax=docker/dockerfile:1.4
FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o kubi8p ./cmd/main.go

# Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/kubi8p /usr/local/bin/kubi8p

ENTRYPOINT ["/usr/local/bin/kubi8p"]