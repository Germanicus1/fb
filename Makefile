.DEFAULT_GOAL := build

BINARY := fb
PKG    := ./cmd/fb

# These api tests call the live Flow Boards API (using ~/.fb/config.yaml) and
# rewrite api/testdata. `make test` skips them; `make test-live` runs them.
LIVE_TESTS := ^(TestAPIFilteringCapabilities|TestBoardBinRelationship|TestInspectAPIResponseForBoardData)$$

.PHONY: build release install test test-live cover vet fmt fmt-check check clean help

build: ## Build ./fb
	go build -o $(BINARY) $(PKG)

release: ## Build a stripped ./fb
	go build -ldflags="-s -w" -o $(BINARY) $(PKG)

install: ## Install fb into $GOPATH/bin
	go install $(PKG)

test: ## Run tests, skipping the live-API ones
	go test -skip '$(LIVE_TESTS)' ./...

test-live: ## Run all tests, including live-API ones; restores api/testdata after
	go test ./... ; status=$$?; git restore api/testdata; exit $$status

cover: ## Test coverage (live-API tests skipped)
	go test -skip '$(LIVE_TESTS)' -cover ./...

vet: ## go vet
	go vet ./...

fmt: ## gofmt all Go files in place
	gofmt -w .

fmt-check: ## Fail if any Go file needs gofmt
	@out=$$(gofmt -l .); if [ -n "$$out" ]; then echo "needs gofmt:"; echo "$$out"; exit 1; fi

check: vet test ## vet + test

clean: ## Remove the built binary
	rm -f $(BINARY)

help: ## List targets
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-10s %s\n", $$1, $$2}'
