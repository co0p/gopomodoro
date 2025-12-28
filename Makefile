.PHONY: build run test clean install uninstall

BINARY_NAME=gopomodoro
BINARY_PATH=bin/$(BINARY_NAME)
INSTALL_PATH=/usr/local/bin/$(BINARY_NAME)

build:
	@mkdir -p bin
	@go build -o $(BINARY_PATH) ./cmd/gopomodoro

run: build
	@$(BINARY_PATH)

test:
	@go test ./...

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_PATH) $(INSTALL_PATH)
	@echo "Installation complete. Run '$(BINARY_NAME)' to start."

uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_PATH)..."
	@sudo rm -f $(INSTALL_PATH)
	@echo "Uninstall complete."

clean:
	@rm -rf bin/
