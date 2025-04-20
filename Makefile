.PHONY: build run test clean docker-build docker-run docker-compose-up docker-compose-down migrate-up migrate-down dev mock coverage swagger swag-init air-install air-dev tools gen-module bench profile security-check

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
	$(GORUN) cmd/migrate/migrate.go

# Rollback a migration
migrate-down:
	@echo "Running migrations down"
	$(GORUN) cmd/migrate/migrate.go -down

# Hot reload for development using air
air-install:
	go install github.com/cosmtrek/air@latest

# Run the application with hot reload
dev: air-install
	air

# Generate mock objects for testing
mock:
	mockery --all --dir=./internal --output=./test/mocks

# Run tests with coverage
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Update Swagger documentation
swagger: swag-init

# Initialize Swagger docs
swag-init:
	swag init -g cmd/main.go -o ./pkg/docs

# Install commonly used development tools
tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

# Generate a new module with basic files
gen-module:
	@read -p "Enter module name: " module_name; \
	mkdir -p modules/$$module_name; \
	touch modules/$$module_name/model.go; \
	touch modules/$$module_name/dto.go; \
	touch modules/$$module_name/repository.go; \
	touch modules/$$module_name/service.go; \
	touch modules/$$module_name/controller.go; \
	echo "Module $$module_name created with basic files"

# Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem ./...

# CPU profiling
profile:
	$(GORUN) -cpuprofile=cpu.prof -memprofile=mem.prof cmd/main.go
	$(GOCMD) tool pprof -http=:8081 cpu.prof

# Check for security issues
security-check:
	gosec -exclude-generated ./...

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
	@echo " - dev: Run with hot reload (using air)"
	@echo " - air-install: Install air for hot reload"
	@echo " - mock: Generate mock objects"
	@echo " - coverage: Run tests and generate coverage report"
	@echo " - swagger: Update Swagger documentation"
	@echo " - swag-init: Initialize Swagger documentation"
	@echo " - tools: Install development tools"
	@echo " - gen-module: Generate a new module with basic files"
	@echo " - bench: Run benchmarks"
	@echo " - profile: Run the application with profiling"
	@echo " - security-check: Check for security issues"
