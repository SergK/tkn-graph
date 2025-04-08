PACKAGE=github.com/sergk/tkn-graph/pkg/cmd

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -v
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
CURRENT_DIR=$(shell pwd)

# Binary name
BINARY_NAME=tkn-graph

# Versioning
GIT_DESCRIBE=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X $(PACKAGE)/version.cliVersion=$(GIT_DESCRIBE)"

override GCFLAGS +=all=-trimpath=${CURRENT_DIR}

.DEFAULT_GOAL:=help

# Directories
SRC_DIR=./cmd/graph
DIST_DIR=./dist

BIN_DIR ?= ${CURRENT_DIR}/bin
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build targets
.PHONY: build
build: ## build the binary
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) -gcflags '${GCFLAGS}' $(SRC_DIR)/...

# Test targets
.PHONY: test
test: ## run tests
	KUBECONFIG=${CURRENT_DIR}/hack/kubeconfig-stub.yaml $(GOTEST) -v -coverprofile=coverage.out ./...

.PHONY: lint
lint: golangci-lint ## Run go lint
	$(GOLANGCI_LINT) run -v -c .golangci.yaml ./...

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix -v -c .golangci.yaml ./...

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
	rm -rf $(DIST_DIR) $(BIN_DIR)

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

GOLANGCI_LINT = ${CURRENT_DIR}/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,v1.64.7)

GITCHGLOG = $(BIN_DIR)/git-chglog
.PHONY: git-chglog
git-chglog: $(BIN_DIR) ## Download git-chglog locally if necessary.
	$(call go-install-tool,$(GITCHGLOG),github.com/git-chglog/git-chglog/cmd/git-chglog,v0.15.4)

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(BIN_DIR) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
