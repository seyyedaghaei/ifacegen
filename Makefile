SHELL := /bin/bash

GOLANGCI_LINT_VERSION ?= v2.11.3
GO ?= go

.PHONY: help fmt test vet lint ci build clean

help:
	@echo "Targets:"
	@echo "  fmt   - gofmt all Go files"
	@echo "  test  - run unit tests"
	@echo "  vet   - run go vet"
	@echo "  lint  - run golangci-lint (installs if missing)"
	@echo "  ci    - fmt check + test + vet + lint"
	@echo "  build - build ifacegen binary into ./ifacegen"
	@echo "  clean - remove build artifacts"

fmt:
	@$(GO) fmt ./...

test:
	@$(GO) test ./...

vet:
	@$(GO) vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..." ; \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) ; \
	}
	@golangci-lint run --timeout=5m

ci:
	@# Fail if gofmt would change files
	@test -z "$$(gofmt -l .)"
	@$(MAKE) test vet lint

build:
	@$(GO) build -o ifacegen ./cmd/ifacegen

clean:
	@rm -f ifacegen

