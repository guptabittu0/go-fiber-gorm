.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down migrate-up migrate-down

# App name
APP_NAME=fiber-gorm-api

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
build:
	$(GOBUILD) -o ./tmp/$(APP_NAME) -v ./cmd/main.go

# Run the application
run:
	$(GORUN) cmd/main.go

# Test the application
test:
	$(GOTEST) -v ./...

# Clean the binary
clean:
	$(GOCLEAN)
	rm -f ./tmp/$(APP_NAME)

# Download dependencies
deps:
	$(GOGET) -u
	$(GOMOD) tidy

# Build the docker image
docker-build:
	docker build -t $(APP_NAME) .

# Run the docker image
docker-run:
	docker run -p 8080:8080 $(APP_NAME)

# Run the application with docker-compose
docker-compose-up:
	docker-compose up -d

# Stop the application with docker-compose
docker-compose-down:
	docker-compose down

# Lint the code
lint:
	golangci-lint run

# Run a migration
migrate-up:
	@echo "Running migrations up"
	$(GORUN) cmd/migrate/up.go

# Rollback a migration
migrate-down:
	@echo "Running migrations down"
	$(GORUN) cmd/migrate/down.go

# Show help
help:
	@echo "Available commands:"
	@echo " - build: Build the application"
	@echo " - run: Run the application"
	@echo " - test: Test the application"
	@echo " - clean: Clean the binary"
	@echo " - deps: Download dependencies"
	@echo " - docker-build: Build the docker image"
	@echo " - docker-run: Run the docker image"
	@echo " - docker-compose-up: Run the application with docker-compose"
	@echo " - docker-compose-down: Stop the application with docker-compose"
	@echo " - lint: Lint the code"
	@echo " - migrate-up: Run migrations up"
	@echo " - migrate-down: Roll back migrations"