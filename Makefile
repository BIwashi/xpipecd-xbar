.DEFAULT_GOAL := help

.PHONY: build/cli
build/cli: BUILD_VERSION ?= $(shell git describe --tags --match "cli-*" --always --abbrev=7 | cut -c9-)
build/cli: BUILD_COMMIT ?= $(shell git rev-parse HEAD)
build/cli: BUILD_REPO ?= github.com/BIwashi/xpipecd-xbar
build/cli: BUILD_DATE ?= $(shell date -u '+%Y%m%d-%H%M%S')
ifdef BUILD_TAGS
	BUILD_TAGS_ARG := -tags $(BUILD_TAGS)
endif
build/cli: BUILD_OPTS ?= -trimpath $(BUILD_TAGS_ARG) -ldflags "$(BUILD_LDFLAGS_PREFIX).version=$(BUILD_VERSION) $(BUILD_LDFLAGS_PREFIX).gitRepository=$(BUILD_REPO) $(BUILD_LDFLAGS_PREFIX).gitCommit=$(BUILD_COMMIT) $(BUILD_LDFLAGS_PREFIX).buildDate=$(BUILD_DATE) -w"
build/cli: BUILD_OS ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f1)
build/cli: BUILD_ARCH ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f2)
build/cli: CGO_ENABLED ?= 0
build/cli: BUILD_ENV ?= GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) CGO_ENABLED=$(CGO_ENABLED)
build/cli: BIN_SUFFIX ?= -$(BUILD_OS)-$(BUILD_ARCH)
build/cli: ## Build cli ## make build/cli
	$(BUILD_ENV) go build $(BUILD_OPTS) -o .artifacts/xpipecd-xbar$(BIN_SUFFIX) ./cmd

.PHONY: setup/cli
setup/cli: build/cli
setup/cli: BUILD_OS ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f1)
setup/cli: BUILD_ARCH ?= $(shell go version | cut -d ' ' -f4 | cut -d/ -f2)
setup/cli: BIN_SUFFIX ?= -$(BUILD_OS)-$(BUILD_ARCH)
setup/cli: t ?= 30s
setup/cli: ## Setup cli ## make setup/cli
	mkdir -p $(HOME)/Library/Application\ Support/xbar/plugins
	mkdir -p $(HOME)/Library/Application\ Support/xbar/plugins/.artifacts
	ln -sf $(CURDIR)/.artifacts/xpipecd-xbar$(BIN_SUFFIX) $(HOME)/Library/Application\ Support/xbar/plugins/.artifacts/xpipecd-xbar
	ln -sf $(CURDIR)/xpipecd-xbar.sh $(HOME)/Library/Application\ Support/xbar/plugins/xpipecd-xbar.$(t).sh

.PHONY: lint/cli
lint/cli: ## Format and lint cli code ## make lint/cli
lint/cli: STRICT_GO_IMPORTS_EXISTS = $(shell which strictgoimports)
lint/cli:
	@if [ ! $(STRICT_GO_IMPORTS_EXISTS) ]; then \
		echo 'strictgoimports is not exists'; \
		go install github.com/momotaro98/strictgoimports/cmd/strictgoimports@latest; \
    fi
	strictgoimports -w -exclude "*.mock.go,*.pb.go" -local "github.com/BIwashi/xpipecd-xbar" .
	golangci-lint run --fix -c .golangci-lint.yml ./...

.PHONY: run/cli
run/cli: ## Run cli ## make run/cli k=${PIPECD_API_KEY} h=pipecd.jp:443
run/cli: k ?= "" # PIPECD_API_KEY
run/cli: h ?= "pipecd.jp:443" # PIPECD_HOST
run/cli:
	go run ./cmd/main.go pipectl --api-key=${k} --host=${h}

##### HELP #####

.PHONY: help
help: ## Display this help screen ## make or make help
	@echo ""
	@echo "Usage: make SUB_COMMAND argument_name=argument_value"
	@echo ""
	@echo "Command list:"
	@echo ""
	@printf "\033[36m%-30s\033[0m %-50s %s\n" "[Sub command]" "[Description]" "[Example]"
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | perl -pe 's%^([/a-zA-Z_-]+):.*?(##)%$$1 $$2%' | awk -F " *?## *?" '{printf "\033[36m%-30s\033[0m %-50s %s\n", $$1, $$2, $$3}'
