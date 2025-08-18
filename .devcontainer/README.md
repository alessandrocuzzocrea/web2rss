# www2rss Development with GitHub Codespaces

This project is configured to work seamlessly with GitHub Codespaces using a devcontainer.

## 🚀 Quick Start

1. **Open in Codespaces**: Click the "Code" button on GitHub and select "Create codespace on main"
2. **Wait for setup**: The devcontainer will automatically install all dependencies
3. **Start developing**: Run `make dev` to start the development server with hot reload

## 📋 Available Commands

```bash
# Development
make dev            # Start development server with hot reload
make build          # Build the application
make test           # Run tests
make lint           # Run linter

# Database
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations
make sqlc-generate  # Generate SQLC code
make setup-db       # Setup database (migrate + generate)

# Utilities
make fmt            # Format code
make clean          # Clean build artifacts
```

## 🔧 Development Environment

The devcontainer includes:

- **Go 1.24** with all standard tools
- **Database tools**: sqlite3, migrate, sqlc
- **Linting**: golangci-lint
- **Hot reload**: air for development
- **VS Code extensions**: Go, GitHub Copilot, and more

## 🌐 Accessing Your Application

- The server runs on port 8080
- In Codespaces, it will be automatically forwarded
- Access URL: `https://CODESPACE_NAME-8080.app.github.dev`

## 📁 Project Structure

```
www2rss/
├── .devcontainer/          # Devcontainer configuration
├── cmd/www2rss/           # Application entry point
├── internal/              # Internal packages
│   ├── app/              # Application logic
│   └── db/               # Database models and queries
├── db/                   # Database migrations and queries
└── data/                 # SQLite database files
```

## 🔄 Hot Reload

The development server uses Air for hot reloading:
- Automatically rebuilds on Go file changes
- Excludes test files and vendor directories
- Logs build errors to `build-errors.log`

## 🗃️ Database

- Uses SQLite for simplicity
- Migrations in `db/migrations/`
- SQLC generates type-safe Go code from SQL
- Database file stored in `data/www2rss.sqlite3`

## 🛠️ Troubleshooting

If something doesn't work:

1. **Rebuild container**: Command Palette → "Codespaces: Rebuild Container"
2. **Check logs**: Look at the setup output during container creation
3. **Manual setup**: Run `.devcontainer/post-create.sh` manually

## 📝 Notes

- The devcontainer automatically sets up the database and generates code
- Git configuration may need to be set manually
- All tools are pre-installed and ready to use
