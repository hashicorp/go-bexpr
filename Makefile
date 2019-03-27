GOTEST_PKGS=$(shell go list ./... | grep -v examples)

TEST_RESULTS?="/tmp/test-results"
grammar.go: grammar.peg
	@echo "Regenerating Parser"
	@go generate ./

generate: grammar.go

test: generate
	@go test $(GOTEST_PKGS)

test-ci: generate
	@gotestsum --junitfile $(TEST_RESULTS)/gotestsum-report.xml -- $(GOTEST_PKGS)

bench: generate
	@go test -bench . $(GOTEST_PKGS)

coverage: generate
	@go test -coverprofile /tmp/coverage.out $(GOTEST_PKGS)
	@go tool cover -html /tmp/coverage.out

fmt: generate
	@gofmt -w -s

examples: simple expr-parse expr-eval filter

simple:
	@go build ./examples/simple

expr-parse:
	@go build ./examples/expr-parse

expr-eval:
	@go build ./examples/expr-eval

filter:
	@go build ./examples/filter

deps:
	@go get github.com/mna/pigeon@master
	@go get golang.org/x/tools/cmd/goimports
	@go get golang.org/x/tools/cmd/cover
	@go mod tidy

.PHONY: generate test coverage fmt deps bench examples expr-parse expr-eval filter

