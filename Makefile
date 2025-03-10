.PHONY: all build run clean docs help vercel vercel-dev dev docker-build docker-run docker-dev docker-stop docker-clean docker-restart

# Default target
.DEFAULT_GOAL := help

# Variables
APP_NAME := golang-template
BUILD_DIR := build
MAIN_FILE := cmd/api/main.go
BINARY_NAME := $(BUILD_DIR)/api
GO_BIN := $(shell go env GOPATH)/bin

# Go build flags
LDFLAGS := -w -s

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  dev          - Run with hot reload (requires Air)"
	@echo "  clean        - Remove build artifacts"
	@echo "  docs         - Generate API documentation"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-dev   - Start Docker Compose for development"
	@echo "  docker-stop  - Stop Docker Compose"
	@echo "  docker-clean - Clean Docker containers and images for this project"
	@echo "  docker-restart - Restart Docker containers"
	@echo "  vercel       - Deploy to Vercel"
	@echo "  vercel-dev   - Run Vercel dev environment"
	@echo "  help         - Show this help message"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_NAME)"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_FILE)

# Run with hot reload
dev:
	@echo "Running $(APP_NAME) with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	elif [ -f $(GO_BIN)/air ]; then \
		$(GO_BIN)/air; \
	else \
		echo "Air not found, installing compatible version..."; \
		go install github.com/cosmtrek/air@v1.40.4; \
		if [ -f $(GO_BIN)/air ]; then \
			$(GO_BIN)/air; \
		else \
			echo "Failed to find air binary after installation."; \
			echo "Using go run with file watching..."; \
			while true; do \
				go run $(MAIN_FILE) & PID=$$!; \
				echo "Server started with PID: $$PID"; \
				inotifywait -e modify -e create -e delete -r --exclude '(\.git|tmp|build)' .; \
				echo "Changes detected, restarting..."; \
				kill -9 $$PID; \
				sleep 1; \
			done; \
		fi; \
	fi


clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) tmp docs/swagger

# Generate API documentation
# Generate API documentation
docs:
	@echo "Generating API documentation..."
	@if [ -f ./scripts/swagger.sh ]; then \
		chmod +x ./scripts/swagger.sh; \
		./scripts/swagger.sh; \
	elif [ -f ./script/swagger.sh ]; then \
		chmod +x ./script/swagger.sh; \
		./script/swagger.sh; \
	else \
		echo "Swagger script not found. Please create it in ./scripts/swagger.sh or ./script/swagger.sh"; \
		exit 1; \
	fi
	@echo "Swagger documentation generated at docs/swagger/"
# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

# Start docker-compose for development with hot reload
docker-up:
	@echo "Starting Docker Compose for development..."
	@docker-compose up -d
	@echo "Application is running at http://localhost:8080"

# Stop docker-compose
docker-down:
	@echo "Stopping Docker Compose..."
	@docker-compose down

# Logs for Docker containers
docker-logs:
	@echo "Showing Docker container logs..."
	@if [ -z "$(shell docker-compose ps -q)" ]; then \
		echo "No running containers found. Start containers first with 'make docker-dev'"; \
	else \
		docker-compose logs -f; \
	fi

# Clean Docker artifacts related to this project
docker-clean:
	@echo "Cleaning Docker artifacts for $(APP_NAME)..."
	@docker-compose down
	@echo "Removing containers..."
	@docker ps -a | grep $(APP_NAME) | awk '{print $$1}' | xargs -r docker rm -f
	@echo "Removing images..."
	@docker images | grep $(APP_NAME) | awk '{print $$3}' | xargs -r docker rmi -f
	@echo "Pruning unused volumes..."
	@docker volume prune -f
	@echo "Docker cleanup complete"

# Restart Docker containers
docker-restart:
	@echo "Restarting Docker containers..."
	@docker-compose down
	@docker-compose up -d
	@echo "Containers restarted. Application is running at http://localhost:8080"

# Vercel deployment
vercel:
	@echo "Deploying to Vercel..."
	@if command -v vercel > /dev/null; then \
		vercel --prod; \
	else \
		echo "Vercel CLI not found. Install with 'npm install -g vercel'"; \
		exit 1; \
	fi


# Run Vercel dev environment
vercel-dev:
	@echo "Running Vercel dev environment..."
	@if command -v vercel > /dev/null; then \
		vercel dev; \
	else \
		echo "Vercel CLI not found. Install with 'npm install -g vercel'"; \
		exit 1; \
	fi

all: clean build
