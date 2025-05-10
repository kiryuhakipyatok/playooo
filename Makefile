-include .env
export

all: build test

build:
	@echo "Building..."
	@go build -o main.exe cmd/app/main.go

run:
	@go run cmd/app/main.go

docker-run:
	@docker compose up --build

docker-down:
	@docker compose down

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -f main

up:
	@goose -dir=$(MIGRATION_PATH) up

down:
	@goose -dir=$(MIGRATION_PATH) down

status:
	@goose -dir=$(MIGRATION_PATH) status

reset:
	@goose -dir=$(MIGRATION_PATH) reset