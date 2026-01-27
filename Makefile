.PHONY: build test install clean

BINARY_NAME=gopomodoro
BUILD_DIR=bin

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gopomodoro

test:
	go test ./...

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
