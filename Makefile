help:
	@echo "deps   -> Get all dependencies"
	@echo "parser -> Generates the sample parser"
	@echo "tests  -> Run all tests"

deps:
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure

parser:
	@pigeon -o "./parser/parser.go" "./parser/nlp.peg"

tests:
	@go test -v -race ./...