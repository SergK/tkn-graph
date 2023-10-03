# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
GOFMT=$(GOCMD) fmt

# Binary name
BINARY_NAME=tkn-graph

# Directories
SRC_DIR=./cmd/graph
BIN_DIR=./bin

# Build targets
.PHONY: build
build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/...

# Test targets
.PHONY: test
test:
	$(GOTEST) -v ./...

# Lint targets
.PHONY: lint
lint:
	$(GOLINT) ./...

# Format targets
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Clean targets
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
