PIGEON_VERSION ?= v1.3.0

help:
	@echo "deps   -> Download all dependencies"
	@echo "parser -> Generates the sample parser"
	@echo "tests  -> Run all tests"

deps:
	@go mod download

parser:
	@go run github.com/mna/pigeon@$(PIGEON_VERSION) -o "./parser/parser.go" "./parser/nlp.peg"

tests:
	@go test -v -race ./...

.PHONY: help deps parser tests
