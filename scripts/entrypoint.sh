#!/bin/sh

set -e

echo "Starting www2rss..."

# Database path
DB_PATH="/root/data/www2rss.sqlite3"
MIGRATIONS_PATH="/root/db/migrations"

# Create data directory
mkdir -p /root/data

# Check if migrations directory exists
if [ -d "$MIGRATIONS_PATH" ]; then
    echo "Running database migrations..."

    # Run migrations using golang-migrate with correct sqlite URL format
    migrate -path "$MIGRATIONS_PATH" -database "sqlite3://$DB_PATH" up

    echo "✅ Migrations completed successfully"
else
    echo "⚠️  No migrations directory found at $MIGRATIONS_PATH"
fi

# Start the application
echo "🚀 Starting www2rss application..."
exec ./www2rss
