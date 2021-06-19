.PHONY: all fmt build test utils doc

GOPATH ?= $(CURDIR)
GOBIN ?= $(GOPATH)/bin

GODOC_PORT ?= 9090

GO ?= go
GODOC ?= $(GOBIN)/godoc
GOGET ?= $(GO) get -v

all: fmt build

fmt:
	$(GO) fmt ./...
	$(GO) mod tidy || true

build:
	$(GO) get -v ./...

test:
	$(GO) test -v ./...

$(GODOC):
	$(GOGET) golang.org/x/tools/cmd/godoc

doc: $(GODOC)
	@echo "http://127.0.0.1:$(GODOC_PORT)"
	@$(GODOC) -http=:$(GODOC_PORT) \
		-index -links=true

utils: $(GODOC)
