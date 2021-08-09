package golang

import "grosf-gh/internal/util"

func createWorkflows() error {
	path := "./" + Name + "/.github"

	// Create folder .github inside the path
	util.CreateFolder(path)

	// Create folder .github/workflows inside the path
	workflowsPath := path + "/workflows"
	util.CreateFolder(workflowsPath)

	// Create  a golang github workflow yaml file inside workflowsPath
	workflowPath := workflowsPath + "/golang.yml"
	workflow := `name: Go
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.16

      - name: Build
        run: go build -v ./...

      - name: Lint and Vet
        run: |
          curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.41.1
          golangci-lint run
          go vet .

      - name: Test
        run: |
          go test ./... -race -covermode=atomic -coverprofile=coverage.out 2>&1 | \
            perl -p -e 'if(/coverage: (\d+.\d)%/) {\
            die "coverage too low: $1" if ($1 < 75) }'

`

	return util.CreateFile(workflowPath, workflow)
}
