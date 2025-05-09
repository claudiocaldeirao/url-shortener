APP_NAME = bootstrap
BUILD_DIR = build
SRC_FILE = app/cmd/main.go

.PHONY: build-UrlShortenerFunction run clean test-docker

## Build the Go Lambda binary
build:
	@echo "🔨 Building Go binary..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

## Run the app locally using SAM
run:
	@echo "🚀 Starting SAM local API Gateway..."
	sam local start-api

## Clean build artifacts
clean:
	@echo "🧹 Cleaning up build directory..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete."

## Test Docker is ready
test-docker:
	@echo "🔎 Checking Docker status..."
	@docker ps > /dev/null 2>&1 && echo "✅ Docker is running!" || echo "❌ Docker is NOT running!"
