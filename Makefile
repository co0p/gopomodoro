.PHONY: build test install clean release run bundle

BINARY_NAME=gopomodoro
BUILD_DIR=bin
RELEASE_OS?=darwin
RELEASE_ARCHS?=$(shell go env GOARCH)

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gopomodoro

test:
	go test ./...

run:
	go run ./cmd/gopomodoro

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)

release: build
	@echo "Building release version(s)..."
	@mkdir -p $(BUILD_DIR)
	@for arch in $(RELEASE_ARCHS); do \
		GOOS=$(RELEASE_OS) GOARCH=$$arch go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(RELEASE_OS)-$$arch ./cmd/gopomodoro; \
	done

bundle:
	@mkdir -p $(BUILD_DIR)/GoPomodoro.app/Contents/MacOS
	@mkdir -p $(BUILD_DIR)/GoPomodoro.app/Contents/Resources
	@cp bundle/macos/GoPomodoro.app/Contents/Info.plist $(BUILD_DIR)/GoPomodoro.app/Contents/Info.plist
	@cp -R bundle/macos/GoPomodoro.app/Contents/Resources $(BUILD_DIR)/GoPomodoro.app/Contents/
	go build -o $(BUILD_DIR)/GoPomodoro.app/Contents/MacOS/$(BINARY_NAME) ./cmd/gopomodoro

clean:
	rm -rf $(BUILD_DIR)
