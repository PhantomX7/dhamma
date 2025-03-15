dep:
	go mod tidy
	go mod vendor

dev:
	go build -o bin/main main.go
	./bin/main

test:
	go test ./... -coverprofile cp.out

test-html:
	go test $(go list ./... | grep -v /mock/) -coverprofile cp.out
	go tool cover -html=cp.out

migrate:
	npx sequelize db:migrate

refresh:
	npx sequelize db:migrate:undo:all
	npx sequelize db:migrate

seed:
	go run ./seeder/main.go

build:
	set GOOS=linux&& set GOARCH=amd64&& go build -o bin/main app/main.go

swag:
	swag init -d app


test-integration:
	go test -v ./test/integration/...