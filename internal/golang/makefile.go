package golang

import "github.com/ezegrosfeld/cli/internal/util"

func createMakefile() error {
	path := "./" + Name + "/Makefile"

	// Create the Makefile content
	content := `.PHONY: test
test:
	@echo "=> Running tests"
	@go test ./... -covermode=atomic -count=1 -race

.PHONY: test-cover
test-cover:
	@echo "=> Running tests and generating report"
	@go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -count=1
	@go tool cover -html=/tmp/coverage.out -o coverage.html

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: build
build:
	@go build main.go

.PHONY: start
start:
	@go run main.go`

	// Create the Makefile
	return util.CreateFile(path, content)

}
