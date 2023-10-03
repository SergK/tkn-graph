# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
GOFMT=$(GOCMD) fmt

# Binary name
BINARY_NAME=tkn-graph

.DEFAULT_GOAL:=help

# Directories
SRC_DIR=./cmd/graph
BIN_DIR=./bin

# Build targets
.PHONY: build
build: ## build the binary
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/...

# Test targets
.PHONY: test
test: ## run tests
	$(GOTEST) -v ./...

# Lint targets
.PHONY: lint
lint: ## run linter
	$(GOLINT) ./...

# Format targets
.PHONY: fmt
fmt: ## run gofmt
	$(GOFMT) ./...

# Clean targets
.PHONY: clean
clean: ## remove the binary
	rm -rf $(BIN_DIR)

# make CI run all targets
.PHONY: all
all: lint test fmt build ## run all targets: lint test fmt build

# Help target
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
