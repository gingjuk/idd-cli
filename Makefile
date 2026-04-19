.PHONY: build run test clean install lint fmt coverage check-coverage

BINARY_NAME=idd-cli
VERSION=1.0.0
GO=go
LINT:=golangci-lint-v2
GOCOVBIN:=$(shell $(GO) env GOPATH)/bin/go-test-coverage

build:
	$(GO) build -o bin/$(BINARY_NAME) ./cmd/idd-cli

run: build
	./bin/$(BINARY_NAME) run .

test:
	CGO_ENABLED=1 $(GO) test -v -race ./...

clean:
	rm -rf bin/
	rm -f idd-report.json
	rm -f coverage.out

install:
	$(GO) install ./cmd/idd-cli

lint:
	$(LINT) config verify
	$(LINT) run ./...
	npx markdownlint-cli2

fmt:
	gofmt -w .
	goimports -w .

coverage:
	$(GO) test -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./...

install-gocov:
	$(GO) install github.com/vladopajic/go-test-coverage/v2@latest

check-coverage: coverage
	$(GOCOVBIN) --config=.testcoverage.yml

.DEFAULT_GOAL := build