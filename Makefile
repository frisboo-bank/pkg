GOPATH:=$(shell go env GOPATH)

MODULE  := frisboo-bank/pkg

# Tool versions
GCI_VERSION := latest
GO_VULN_CHECK_VERSION := latest
GOFUMPT_VERSION := latest
GOLANGCI_VERSION := latest
GOLINES_VERSION := latest
MOCKERY_VERSION  := latest
REVIVE_VERSION := latest
STATIC_CHECK_VERSION := latest

#
# default target
#
.PHONY: all
all: help

#
# QUALITY
#
## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v
	go run github.com/segmentio/golines@$(GOLINES_VERSION) -m 120 -w --ignore-generated .
	go run github.com/daixiang0/gci@$(GCI_VERSION) write --skip-generated -s standard -s "prefix($(MODULE))" -s default -s blank -s dot --custom-order  .
	go run mvdan.cc/gofumpt@$(GOFUMPT_VERSION) -l -w .

## audit: Run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@$(STATIC_CHECK_VERSION) -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@$(GO_VULN_CHECK_VERSION) ./...
	go test -race -buildvcs -vet=off ./...

## lint: Run linters
.PHONY: lint
lint:
	go run github.com/mgechev/revive@$(REVIVE_VERSION) -config revive-config.toml -formatter friendly ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION) run ./...

.PHONY: update
update:
	@go get -u

.PHONY: deps-reset
deps-reset:
	git checkout -- go.mod
	go mod tidy

.PHONY: deps-upgrade
deps-upgrade:
	go get -u -t -v ./...
	go mod tidy

.PHONY: deps-cleancache
deps-cleancache:
	go clean -modcache

#
# HELPERS
#
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /' | sort
