.PHONY: all check test lint fmt bench

GO_FILES:=$(shell find -type f -iname '*.go')

all: check

check: lint test

test: $(GO_FILES)
	go test -v .

fmt: $(GO_FILES)
	gofmt -w .

lint: $(GO_FILES)
	golangci-lint run .

bench:
	go test -bench .