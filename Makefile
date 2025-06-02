GOPATH:=$(shell go env GOPATH)

.PHONY: lint
lint:
	revive -config revive-config.toml -formatter friendly ./...
	staticcheck ./...
	golangci-lint run ./...

.PHONY: format
format:
	golines -m 120 -w --ignore-generated .
	# gci write --skip-generated -s standard -s "prefix(github.com/mehdihadeli/go-food-delivery-microservices)" -s default -s blank -s dot --custom-order  .
	gofumpt -l -w .

.PHONY: update
update:
	@go get -u

.PHONY: tidy
tidy:
	go mod tidy

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
