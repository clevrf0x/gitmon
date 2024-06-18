BUILD_DIR := bin
BUILD_OUTPUT := $(BUILD_DIR)/gitmon
SRC_DIR := cmd/gitmon
SRC_FILE := $(SRC_DIR)/main.go

.DEFAULT_GOAL := build

.PHONY: all test build run clean

all: test build

test:
	go test -v ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_OUTPUT) $(SRC_FILE)

run: build
	$(BUILD_OUTPUT) $(ARGS)

clean:
	rm -rf $(BUILD_DIR)

# Ensure the build directory exists before building
$(BUILD_OUTPUT): | $(BUILD_DIR)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

