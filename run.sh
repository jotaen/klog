#!/bin/sh

# Install all dependencies
run_install() {
	go get -t ./...
}

# Compile to ./out/klog
# Takes two positional arguments:
# - The version (e.g.: v1.2)
# - The build hash (7 chars hex)
run_build() {
	go build \
	  -ldflags "-X 'main.BinaryVersion=$1' -X 'main.BinaryBuildHash=$2'" \
	  -o ./out/klog \
	  klog.go
}

# Execute all tests
run_test() {
	go test ./...
}

# Reformat all code
run_format() {
	go fmt ./...
}

# Static code (style) analysis
run_lint() {
  set -o errexit
  go vet ./...
  staticcheck ./...
}

# Run CLI from sources “on the fly”
# Passes through all input args
run_cli() {
	go run klog.go "$@"
}
