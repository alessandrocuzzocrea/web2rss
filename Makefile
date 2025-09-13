# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=web2rss
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKER_IMAGE=web2rss:latest

.PHONY: all build clean test coverage deps docker docker-run help dev lint migrate-up migrate-down sqlc-generate dump-schema

all: deps test build ## Run deps, test and build

dev: ## Start development server with hot reload
	air -c .air.toml

build: ## Build the binary file
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/web2rss

clean: ## Remove binary and test cache
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out coverage.html

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

coverage: test ## Run tests and show coverage
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

# Cross compilation
build-linux: ## Build for Linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/web2rss

# Docker commands
docker: ## Build docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker ## Run docker container
	docker run -p 8080:8080 $(DOCKER_IMAGE)

docker-compose-up: ## Run with docker-compose
	docker-compose up --build

docker-compose-down: ## Stop docker-compose
	docker-compose down

# Development commands
dev: ## Run in development mode
	$(GOCMD) run ./cmd/web2rss

watch: ## Run with file watcher (requires entr)
	find . -name "*.go" | entr -r $(GOCMD) run ./cmd/web2rss

# Linting and formatting
fmt: ## Format code
	$(GOCMD) fmt ./...

lint: ## Run golangci-lint
	golangci-lint run

# Database commands
migrate-up: ## Run database migrations up
	migrate -database "sqlite://data/web2rss.sqlite3" -path db/migrations up

migrate-down: ## Run database migrations down
	migrate -database "sqlite://data/web2rss.sqlite3" -path db/migrations down

sqlc-generate: ## Generate SQLC code
	sqlc generate

setup-db: migrate-up sqlc-generate ## Setup database (migrate + generate code)

# Help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

dump-schema: ## Dump current database schema to db/schema.sql
	sqlite3 ./data/web2rss.sqlite3 ".schema" > ./db/schema.sql
