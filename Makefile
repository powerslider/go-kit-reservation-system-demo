all: clean test doc build

doc:
	@echo ">>> Generate Swagger API Documentation..."
	swag init --generalInfo cmd/reservation/main.go

build:
	@echo ">>> Building Application..."
	go build -o bin/reservations cmd/reservation/main.go

test:
	@echo ">>> Running Unit Tests..."
	go test -race ./...

cover-test:
	@echo ">>> Running Tests with Coverage..."
	go test -race ./... -coverprofile=coverage.txt -covermode=atomic

clean:
	@echo ">>> Removing binaries..."
	@rm -rf bin/*
