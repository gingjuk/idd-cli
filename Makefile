.PHONY: build run test clean install

BINARY_NAME=idd-cli
VERSION=1.0.0
GO=go

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
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

.DEFAULT_GOAL := build