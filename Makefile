.PHONY: build run docker-up docker-down docker-rebuild clean

build:
	go build -o bin/main ./cmd/api

run:
	go run ./cmd/api/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose down
	docker-compose up --build

clean:
	rm -rf bin/
	docker-compose down -v
