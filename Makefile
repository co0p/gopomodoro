.PHONY: build run test clean

BINARY_NAME=gopomodoro
BINARY_PATH=bin/$(BINARY_NAME)

build:
	@mkdir -p bin
	@go build -o $(BINARY_PATH) ./cmd/gopomodoro

run: build
	@$(BINARY_PATH)

test:
	@go test ./...

clean:
	@rm -rf bin/
