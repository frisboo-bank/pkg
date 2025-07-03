GOFLAGS ?=
MAKEFLAGS += --no-builtin-rules
.DEFAULT_GOAL := help
.EXPORT_ALL_VARIABLES:

MIN_GO_VERSION := 1.20

GOPATH := $(shell go env GOPATH)

MODULE  := frisboo-bank/pkg

# Tool versions
GCI_VERSION := latest
GO_ENUMS_VERSION := latest
GO_VULN_CHECK_VERSION := latest
GOFUMPT_VERSION := latest
GOIMPORTS_VERSION := latest
GOLANGCI_VERSION := latest
GOLINES_VERSION := latest
MOCKERY_VERSION  := latest
PKGSITE_VERSION := latest
REVIVE_VERSION := latest
STATIC_CHECK_VERSION := latest

#
# DEFAULT TARGET: show help
#
.PHONY: all
all: help

#
# DEVELOPMENT
#
## generate: Run go generate for all packages
.PHONY: generate
generate:
	go generate ./...

## generate-mock: Generate mocks with mockery
.PHONY: generate-mock
generate-mock:
	go run github.com/vektra/mockery/v2@$(MOCKERY_VERSION) --all --output=./mocks

## doc: Generate the documentation
.PHONY: doc
doc:
	go doc ./...

## doc-serve: Serve documentation locally for browsing
.PHONY: doc-serve
doc-serve:
	 go run golang.org/x/pkgsite/cmd/pkgsite@$(PKGSITE_VERSION) -http=localhost:6060

#
# Testing
#
## test: Run all tests
.PHONY: test
test:
	go test -v ./...

## quick-test: Run only quick Go tests
.PHONY: quick-test
quick-test:
	go test -short ./...

## coverage: Run tests with code coverage
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

## coverage-html: Open HTML coverage report in browser
.PHONY: coverage-html
coverage-html:
	go tool cover -html=coverage.out

## bench: Run Go benchmarks
.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

## benchcmp: Compare two benchmark outputs (requires benchcmp)
.PHONY: benchcmp
benchcmp:
	benchcmp old.txt new.txt

## debug: Run tests with Delve debugger
.PHONY: debug
debug:
	dlv test ./...

#
# QUALITY
#
## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy
	go mod verify
	go run golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION) -w .
	go run github.com/segmentio/golines@$(GOLINES_VERSION) -m 120 -w --ignore-generated .
	go run github.com/daixiang0/gci@$(GCI_VERSION) write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order  .
	go run mvdan.cc/gofumpt@$(GOFUMPT_VERSION) -l -w .

## lint: Run linters
.PHONY: lint
lint:
	go run github.com/mgechev/revive@$(REVIVE_VERSION) -config revive-config.toml -formatter friendly ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION) run ./...

## audit: Run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@$(STATIC_CHECK_VERSION) -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@$(GO_VULN_CHECK_VERSION) ./...
	go test -race ./...

## vet: Run go vet on all packages
.PHONY: vet
vet:
	go vet ./...

#
# Dependency Management
#
## deps-reset: Reset go.mod changes
.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod go.sum
	go mod tidy

## deps-upgrade: Upgrade all dependencies
.PHONY: deps-upgrade
deps-upgrade:
	go get -u -t -v ./...
	go mod tidy

## deps-cleancache: Clean Go modules cache
.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache

## deps-tidy: Clean up go.mod and go.sum
.PHONY: deps-tidy
deps-tidy:
	go mod tidy
	go mod verify

## deps-outdated: List outdated dependencies
.PHONY: deps-outdated
deps-outdated:
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all

## modgraph: Print the module dependency graph
.PHONY: modgraph
modgraph:
	go mod graph | sort

#
# MISC
#
## check: Run all checks and tests
.PHONY: check
check: check-go-version tidy lint audit test coverage

## ci: Run all checks and tests (what CI would run)
.PHONY: ci
ci: check

## precommit: Run essential checks before committing
.PHONY: precommit
precommit: check

## check-go-version: Ensure Go version is >= MIN_GO_VERSION
.PHONY: check-go-version
check-go-version:
	@current=$$(go version | awk '{print $$3}' | sed 's/go//'); \
	required=$(MIN_GO_VERSION); \
	lowest=$$(printf '%s\n' "$$required" "$$current" | sort -V | head -n1); \
	if [ "$$lowest" != "$$required" ]; then \
		echo "Go version is too old: found $$current, require Go >= $$required"; exit 1; \
	fi

## clean: Clean build/test artifacts
.PHONY: clean
clean:
	rm -f coverage.out
	rm -rf ./mocks

## tools: Install all required Go tools for development
.PHONY: tools
tools:
	# Linters & formatters
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go install github.com/mgechev/revive@$(REVIVE_VERSION)
	go install github.com/segmentio/golines@$(GOLINES_VERSION)
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	go install github.com/daixiang0/gci@$(GCI_VERSION)
	# Mocks & codegen
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	go install github.com/zarldev/goenums@$(GO_ENUMS_VERSION)
	# Security & static analysis
	go install golang.org/x/vuln/cmd/govulncheck@$(GO_VULN_CHECK_VERSION)
	go install honnef.co/go/tools/cmd/staticcheck@$(STATIC_CHECK_VERSION)
	# Documentation
	go install golang.org/x/pkgsite/cmd/pkgsite@$(PKGSITE_VERSION)

## version: Print project version
.PHONY: version
version:
	@git describe --tags --always || echo "No version tag"

#
# HELPERS
#
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /' | sort
	@echo ""
	@echo "Current Go version: $$(go version | awk '{print $$3}' | sed 's/go//')"

## list: List all available make targets
.PHONY: list
list:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | cut -d':' -f1 | sort
