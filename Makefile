# Go parameters
GO := go
BINARY_NAME := $(shell basename $(CURDIR))
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DATE := $(shell date '+%Y-%m-%d.%H%M%S')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')

# Compiler flags
LD_FLAGS := -X 'main.version=$(VERSION)' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.commitHash=$(COMMIT_HASH)' -X 'main.branch=$(BRANCH)'

# Tool arguments
TAGS := json,yaml,xml

# Default goal
.DEFAULT_GOAL := build

# Targets
.PHONY: all build clean deps install-tools run tags

all: build

build: | bin
	@$(GO) build -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME) cmd/cli/main.go

clean:
	@$(GO) clean
	@rm -rf bin

deps:
	@export GOPRIVATE=github.com/bengrewell && $(GO) get -u ./...

install-tools:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go
	@$(GO) get github.com/fatih/gomodifytags

run:
	@$(GO) run cmd/main.go

tags:
	@gomodifytags -file $(FILE) -all -add-tags $(TAGS) -w

bin:
	@mkdir -p bin
