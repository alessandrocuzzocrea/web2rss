# web2rss

`web2rss` is a lightweight Go application that turns any website into an RSS feed using CSS selectors.

## Features

- **Visual Preview**: Test your selectors in real-time before creating a feed.
- **Scheduled Refresh**: Automatically polls websites and updates your items.
- **Configuration**: Easy setup via environment variables.
- **SQLite Backend**: Fast and portable data storage.

## Getting Started

### Prerequisites

- Go 1.24+
- `sqlc` (optional, for DB code generation)
- `migrate` (optional, for DB migrations)

### Configuration

You can configure the application using environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Path to the SQLite database (default: `./data/web2rss.sqlite3`)
- `DATA_DIR`: Directory for data storage (default: `./data`)

### Run with Makefile

```bash
# Setup database
make setup-db

# Run in development mode
make dev
```

### Direct Run

```bash
# Build the binary
go build ./cmd/web2rss

# Run the app
./web2rss
```

## Database Migrations

To create a new migration:
```bash
migrate create -ext sql -dir db/migrations -seq create_xxxx_table
```

To apply migrations manually:
```bash
migrate -database "sqlite3://data/web2rss.sqlite3" -path db/migrations up
```

## SQLC Generation

To regenerate database code from SQL queries:
```bash
go tool sqlc generate
```
