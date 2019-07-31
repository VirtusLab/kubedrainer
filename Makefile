# Setup variables for the Makefile
NAME := kubedrainer
PKG := github.com/VirtusLab/$(NAME)

# Set POSIX sh for maximum interoperability
SHELL := /bin/sh
PATH  := $(GOPATH)/bin:$(PATH)

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

# Set the main.go path for go command
BUILD_PATH := ./cmd/$(NAME)

# Set any default go build tags
BUILDTAGS :=

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := ${PREFIX}/cross

ARGS ?= $(EXTRA_ARGS)

.DEFAULT_GOAL := help

.PHONY: mod
mod: ## Populates the vendor directory with dependencies
	@echo "+ $@"
	go mod tidy
	go mod vendor

.PHONY: build
build: $(NAME) ## Builds a dynamic executable or package
	@echo "+ $@"

$(NAME): $(wildcard *.go) $(wildcard */*.go)
	@echo "+ $@"
	go build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME) $(BUILD_PATH)

.PHONY: static
static: ## Builds a static executable
	@echo "+ $@"
	CGO_ENABLED=0 go build \
				-tags "$(BUILDTAGS) static_build" \
				${GO_LDFLAGS_STATIC} -o $(NAME) $(BUILD_PATH)

.PHONY: run
run: ## Run the executable, you can use EXTRA_ARGS
	@echo "+ $@"
	go run -tags "$(BUILDTAGS)" ${GO_LDFLAGS} $(BUILD_PATH)/main.go $(ARGS)

.PHONY: clean
clean: ## Cleanup any build binaries or packages
	@echo "+ $@"
	go mod tidy
	$(RM) $(NAME) || echo "Couldn't delete, not there."
	$(RM) test$(NAME) || echo "Couldn't delete, not there."
	$(RM) -r $(BUILDDIR) || echo "Couldn't delete, not there."
	$(RM) coverage.txt || echo "Couldn't delete, not there."

.PHONY: help
help:
	@grep -Eh '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
