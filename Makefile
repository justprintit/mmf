.PHONY: all fmt build test

GOPATH ?= $(CURDIR)
GOBIN ?= $(GOPATH)/bin

GO ?= go
GOGET ?= $(GO) get -v
GOFMT ?= gofmt
GOFMT_FLAGS = -w -l -s

all: fmt build

fmt:
	@find . -name '*.go' | xargs -r $(GOFMT) $(GOFMT_FLAGS)
	$(GO) mod tidy || true

build:
	$(GO) get -v ./...

test:
	$(GO) test -v ./...
