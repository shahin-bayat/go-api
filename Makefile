APP_NAME = go_api
BUILD_DIR = ./bin
MAIN_FILE = ./cmd/api-server/main.go

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

run: build
	@$(BUILD_DIR)/$(APP_NAME)

test:
	@go test -v ./...

# Add a seed target
seed: build
	@$(BUILD_DIR)/$(APP_NAME) --seed
