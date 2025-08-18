#!/bin/bash

echo "🚀 Setting up www2rss development environment..."

# Update package list
sudo apt-get update

# Install additional tools
echo "📦 Installing development tools..."
sudo apt-get install -y \
    sqlite3 \
    curl \
    wget \
    jq \
    tree \
    htop

# Install Go tools
echo "🔧 Installing Go tools..."
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest

# Create data directory
echo "📁 Creating data directory..."
mkdir -p ./data

# Install Go dependencies
echo "📚 Installing Go dependencies..."
go mod download

# Run database migrations if they exist
if [ -d "./db/migrations" ]; then
    echo "🗃️ Setting up database..."
    migrate -database "sqlite://data/www2rss.sqlite3" -path db/migrations up || echo "⚠️ Migration failed or no migrations to run"
fi

# Generate SQLC code
if [ -f "./sqlc.yaml" ]; then
    echo "⚙️ Generating SQLC code..."
    sqlc generate || echo "⚠️ SQLC generation failed"
fi

# Set up Git if not already configured
if [ -z "$(git config --global user.name)" ]; then
    echo "🔧 Setting up Git configuration..."
    echo "Please run:"
    echo "  git config --global user.name 'Your Name'"
    echo "  git config --global user.email 'your.email@example.com'"
fi

echo "✅ Development environment setup complete!"
echo ""
echo "🚀 Quick start commands:"
echo "  make dev     # Start development server with hot reload"
echo "  make build   # Build the application"
echo "  make test    # Run tests"
echo "  make lint    # Run linter"
echo ""
echo "📖 The server will be available on port 8080"
echo "🔗 Access it at: https://CODESPACE_NAME-8080.app.github.dev"
