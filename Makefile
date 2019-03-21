./bexpr/grammar.go: ./bexpr/grammar.peg
	@echo "Regenerating Parser"
	@go generate ./bexpr

generate: ./bexpr/grammar.go

test: generate
	@go test ./bexpr/...

coverage: generate
	@go test -coverprofile /tmp/coverage.out ./bexpr/...
	@go tool cover -html /tmp/coverage.out

fmt: generate
	@gofmt -w -s

deps:
	@go get github.com/mna/pigeon@master
	@go get golang.org/x/tools/cmd/goimports
	@go get golang.org/x/tools/cmd/cover
	@go mod tidy

.PHONY: generate test coverage fmt deps

