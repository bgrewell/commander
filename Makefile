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
.PHONY: all build clean deps install-tools run tags dist

all: build

build: | bin
	@$(GO) build -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME) cmd/cli/main.go

clean:
	@$(GO) clean
	@rm -rf bin

deps:
	@export GOPRIVATE=github.com/bgrewell && $(GO) get -u ./...

install-tools:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go
	@$(GO) get github.com/fatih/gomodifytags

run:
	@$(GO) run cmd/main.go

tags:
	@gomodifytags -file $(FILE) -all -add-tags $(TAGS) -w

bin:
	@mkdir -p bin

dist: dist-windows-x86 dist-windows-amd64 dist-windows-arm dist-linux-amd64 dist-linux-arm dist-osx-amd64 dist-osx-arm

dist-%:
	GOOS=$(word 1, $(subst -, ,$*)) GOARCH=$(word 2, $(subst -, ,$*)) $(GO) build -ldflags "$(LD_FLAGS)" -o dist/$(BINARY_NAME)-$* cmd/cli/main.go

dist-windows-x86:
	GOOS=windows GOARCH=386 $(GO) build -ldflags "$(LD_FLAGS)" -o dist/windows-x86 cmd/cli/main.go

dist-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LD_FLAGS)" -o dist/windows-amd64 cmd/cli/main.go

dist-windows-arm:
	GOOS=windows GOARCH=arm $(GO) build -ldflags "$(LD_FLAGS)" -o dist/windows-arm cmd/cli/main.go

dist-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LD_FLAGS)" -o dist/linux-amd64 cmd/cli/main.go

dist-linux-arm:
	GOOS=linux GOARCH=arm $(GO) build -ldflags "$(LD_FLAGS)" -o dist/linux-arm cmd/cli/main.go

dist-osx-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LD_FLAGS)" -o dist/osx-amd64 cmd/cli/main.go

dist-osx-arm:
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LD_FLAGS)" -o dist/osx-arm cmd/cli/main.go