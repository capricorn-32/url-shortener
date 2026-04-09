FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortener ./cmd/server

FROM alpine:3.20
WORKDIR /app

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /url-shortener /usr/local/bin/url-shortener

EXPOSE 5000
USER app
ENTRYPOINT ["/usr/local/bin/url-shortener"]
