BINARY_NAME=admin-user-service
PORT=8080
DOCKER_IMAGE=admin-user-service
DOCKER_TAG=latest

GOCMD=go

.PHONY: run build deps migrate-up migrate-down docker-build help

run:
	$(GOCMD) run main.go server

build:
	@mkdir -p bin
	$(GOCMD) build -o bin/$(BINARY_NAME) .

deps:
	$(GOCMD) mod download
	$(GOCMD) mod tidy

migrate-up:
	migrate -path migrations -database "mysql://$${DATABASE_USERNAME}:$${DATABASE_PASSWORD}@tcp($${DATABASE_HOST}:$${DATABASE_PORT})/$${DATABASE_DBNAME}" up

migrate-down:
	migrate -path migrations -database "mysql://$${DATABASE_USERNAME}:$${DATABASE_PASSWORD}@tcp($${DATABASE_HOST}:$${DATABASE_PORT})/$${DATABASE_DBNAME}" down 1

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

help:
	@echo "admin-user-service: make run | build | migrate-up | docker-build"
