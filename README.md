# URL Shortener (Go + Redis)

A small URL shortener API built with Gin and Redis.

## Restructured project layout

- `cmd/server`: application entrypoint
- `internal/config`: environment configuration
- `internal/server`: HTTP router setup
- `internal/handler`: request handlers
- `internal/shortener`: short URL generation logic
- `internal/store`: Redis persistence layer

## Run locally

Prerequisites:
- Go 1.22+
- Redis running on `localhost:6379`

Commands:

```bash
go run ./cmd/server
```

## Run with Docker Compose

```bash
docker compose up --build
```

API is available at `http://localhost:5000`.

## API endpoints

- `GET /`: health response
- `POST /create-short-url`: create a short URL
- `GET /:shortUrl`: redirect to the original URL

Sample payload:

```json
{
  "long_url": "https://example.com/very/long/path",
  "user_id": "user-123"
}
```

## Tests

```bash
go test ./...
```
