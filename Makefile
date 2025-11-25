.PHONY: build run test clean docker-up docker-down

BINARY_NAME=pr-reviewer

build:
	go build -o bin/$(BINARY_NAME) ./cmd/server/main.go

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./tests/...

clean:
	go clean
	rm -rf bin/

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down