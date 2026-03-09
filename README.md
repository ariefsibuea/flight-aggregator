# Flight Aggregator

Flight Aggregator is a Go-based flight search service that fans out requests to multiple airline providers concurrently, normalizes response formats into a single unified schema, and returns results through a single REST endpoint.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose v2+
- Go 1.25+

## Setup & Run

### Docker

```bash
git clone https://github.com/ariefsibuea/flight-aggregator.git
cd flight-aggregator
cp .env.example .env
docker compose up --build
# stop the service
docker compose down
```

The API starts on port `8080` and Redis on port `6379`. Verify:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### Local Development

```bash
# Start Redis
docker compose up redis -d

# Run the application
cp .env.example .env
go run cmd/api/main.go
```

## Documentation

- **API Reference:** [`docs/api.md`](docs/api.md)
- **Design:** [`docs/design.md`](docs/design.md)
