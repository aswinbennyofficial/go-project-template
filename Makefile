.PHONY: build run test migrate

build:
    go build -o myapp

run: build
    ./myapp

test:
    go test ./...

migrate:
    goose -dir ./src/db/migrations postgres "user=user password=password dbname=myapp sslmode=disable" up
