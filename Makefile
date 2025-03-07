# Variables
APP_NAME := myapp
VERSION := latest
DOCKER_IMAGE := $(APP_NAME):$(VERSION)
HOST_PORT := 8080
CONTAINER_PORT := 8080
DOCKER_NETWORK := app-network
CPU := 1
MEMORY := 1024m

# Default target
.PHONY: all
all: build

# Build the Docker image
.PHONY: build
build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .

# Run the Docker container
.PHONY: run
run:
	@echo "Running Docker container from image $(DOCKER_IMAGE)..."
	docker run -it -p $(HOST_PORT):$(CONTAINER_PORT) --network $(DOCKER_NETWORK) --memory=$(MEMORY) --cpus=$(CPU) --name $(APP_NAME) $(DOCKER_IMAGE)

# Stop the Docker container
.PHONY: stop
stop:
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# Restart the Docker container
.PHONY: restart
restart: stop run

# Clean Docker resources
.PHONY: clean
clean: stop
	@echo "Removing Docker image $(DOCKER_IMAGE)..."
	docker rmi $(DOCKER_IMAGE) || true

# Show logs
.PHONY: logs
logs:
	docker logs -f $(APP_NAME)

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build    - Build Docker image"
	@echo "  make run      - Run Docker container"
	@echo "  make stop     - Stop Docker container"
	@echo "  make restart  - Restart Docker container"
	@echo "  make clean    - Remove Docker container and image"
	@echo "  make logs     - Show container logs"
	@echo "  make help     - Show this help message"
