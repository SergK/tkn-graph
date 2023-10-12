# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
CURRENT_DIR=$(shell pwd)

# Binary name
BINARY_NAME=tkn-graph

.DEFAULT_GOAL:=help

# Directories
SRC_DIR=./cmd/graph
DIST_DIR=./dist
BIN_DIR=./bin

# Build targets
.PHONY: build
build: ## build the binary
	$(GOBUILD) -o $(DIST_DIR)/$(BINARY_NAME) $(SRC_DIR)/...

# Test targets
.PHONY: test
test: ## run tests
	KUBECONFIG=${CURRENT_DIR}/hack/kubeconfig-stub.yaml $(GOTEST) -v -coverprofile=coverage.out ./...

# Lint targets
.PHONY: lint
lint: ## run linter
	$(GOLINT) ./...

# Format targets
.PHONY: fmt
fmt: ## run gofmt
	$(GOFMT) ./...

.PHONY: vet
vet:  ## Run go vet
	$(GOVET) ./...

# Clean targets
.PHONY: clean
clean: ## remove the binary
	rm -rf $(DIST_DIR)

# make CI run all targets
.PHONY: all
all: lint test fmt vet build ## run all targets: lint test fmt build

# use https://github.com/git-chglog/git-chglog/
.PHONY: changelog
changelog: git-chglog	## generate changelog
ifneq (${NEXT_RELEASE_TAG},)
	$(GITCHGLOG) --next-tag v${NEXT_RELEASE_TAG} -o CHANGELOG.md v0.1.0..
else
	$(GITCHGLOG) -o CHANGELOG.md v0.1.0..
endif

# Help target
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

GITCHGLOG = ${BIN_DIR}/git-chglog
.PHONY: git-chglog
git-chglog: ## Download git-chglog locally if necessary.
	$(call go-get-tool,$(GITCHGLOG),github.com/git-chglog/git-chglog/cmd/git-chglog,v0.15.4)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get -d $(2)@$(3) ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef
