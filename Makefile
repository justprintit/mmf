.PHONY: all fmt build test

GO ?= go

all: fmt build

fmt:
	$(GO) fmt ./...
	$(GO) mod tidy || true

build:
	$(GO) get -v ./...

test:
	$(GO) test -v ./...
