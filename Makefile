SHELL := /bin/bash

REQUIRED_BINS := go
.DEFAULT_TARGET: help

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

setup: ## Check all required tools and install dependencies
	$(if $(shell cat .vscode/.setup 2> /dev/null),$(echo Setup already done!))
	@echo "Verifying required tools... (go)"
	$(foreach bin,$(REQUIRED_BINS),\
    	$(if $(shell command -v $(bin) 2> /dev/null),,$(error Please install `$(bin)`)))

	@echo "Setting up your project..."
	@echo "Installing dependencies..."

	@go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
	@go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -v github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@go install -v golang.org/x/tools/cmd/goimports@latest

	@mkdir -p .vscode
	@touch .vscode/.setup
	@echo "Setup done!"


lint: ## Execute syntatic analysis in the code and autofix minor problems
	@golangci-lint run --fix


version: ## Show the current version of the application
	@echo "Version: $$(git describe --tags --abbrev=0 | tr -d '\n')" 
	

build: ## Build the application
	@.github/scripts/build.sh

install: ## Build and install the application
	@.github/scripts/build.sh install
reinstall: install ## Remove and install the application
	@gs-dev install --uninstall
	@gs-dev install
	
.PHONY: release
release: ## Generate new release on github
	@python3 .github/scripts/new_release.py

.PHONY: docs
docs: ## Generate documentation for the application
	@go get github.com/spf13/cobra/doc@latest
	@go run cmd/docs/docs.go --docs=./docs