# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
GOFMT=$(GOCMD) fmt

# Binary name
BINARY_NAME=myapp

# Directories
SRC_DIR=.
BIN_DIR=bin

# Build targets
.PHONY: build
build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/...

# Test targets
.PHONY: test
test:
	$(GOTEST) -v $(SRC_DIR)/...

# Lint targets
.PHONY: lint
lint:
	$(GOLINT) $(SRC_DIR)/...

# Format targets
.PHONY: fmt
fmt:
	$(GOFMT) $(SRC_DIR)/...
