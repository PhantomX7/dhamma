# Commands
dep: 
	go mod tidy
	go mod vendor

run: 
	go run main.go

build: 
	go build -o bin/main main.go

run-build: build
	./bin/main

test:
	go test -v ./tests

init-docker:
	docker compose up -d --build

up: 
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f