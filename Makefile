.PHONY: build run test clean install

BINARY_NAME=idd-cli
VERSION=1.0.0
GO=go
LINT:=golangci-lint-v2

build:
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/validator

run: build
	./bin/$(BINARY_NAME) run . -v

test:
	$(GO) test -v -race ./...

clean:
	rm -rf bin/
	rm -f idd-report.json

install: build
	install -m 755 bin/$(BINARY_NAME) /usr/local/bin/

lint:
	$(LINT) config verify
	$(LINT) run ./...

fmt:
	gofmt -w .
	goimports -w .

.DEFAULT_GOAL := build