#!/bin/bash

set -e

# Generate test data.
FILE="$(mktemp)"
go run benchmark.go 10000 > "${FILE}"

# Ensure binary is all set.
klog > /dev/null

# Run benchmark.
time klog total --no-warn "${FILE}"
