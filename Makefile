BINARY_NAME		= csscgen
GIT_COMMIT_HASH = $(shell git rev-parse HEAD)
GIT_TREE_STATE 	= $(shell test -n "`git status --porcelain`" && echo "modified" || echo "unmodified")
GIT_TAG     	= $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)

LDFLAGS = -w
LDFLAGS += -X $(GO_PKG)/internal/version.GitCommitHash=$(GIT_COMMIT_HASH)
LDFLAGS += -X $(GO_PKG)/internal/version.GitTreeState=$(GIT_TREE_STATE)
LDFLAGS += -X $(GO_PKG)/internal/version.GitTag=$(GIT_TAG)

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: build-cli
build-cli: 
	go build --ldflags="$(LDFLAGS)" -cover \
	-coverpkg=github.com/akashsinghal/csscgen/cmd/... \
	-o ./bin/${BINARY_NAME} ./cmd/${BINARY_NAME}