#!/bin/bash
# Use this script with https://run.jotaen.net

# Install all dependencies
run::install() {
	go get -t ./...
	go mod tidy
}

# Compile to ./out/klog
# Takes three positional arguments (see `klog.go` file):
# - Version string (e.g.: v1.2)
# - Build hash (7 chars hex)
# - Build date (formatted YYYY-MM-DD)
run::build() {
	go build \
	  -ldflags "\
	    -s -w \
	    -X 'main.BinaryVersion=$1' \
	    -X 'main.BinaryBuildHash=$2' \
	    -X 'main.BinaryBuildDate=$3' \
	    " \
	  -o ./out/klog \
	  klog.go
}

# Execute all tests
run::test() {
	go test ./...
}

# Reformat all code
run::format() {
	go fmt ./...
}

# Static code (style) analysis
run::lint() {
  set -o errexit
  go vet ./...
  staticcheck ./...
}

# Run CLI from sources “on the fly”
# Passes through all input args
run::cli() {
	go run klog.go "$@"
}
