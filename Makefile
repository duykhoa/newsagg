.PHONY: all tidy test run build sqlc

all: test build

tidy:
	go mod tidy

test:
	go test -v ./...

build:
	go build -o bin/newsagg cmd/newsagg/main.go

run:
	go run cmd/newsagg/main.go

sqlc:
	sqlc generate
