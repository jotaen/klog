#!/bin/bash

set -e

FILE="$(mktemp)"
go run benchmark.go 10000 > "${FILE}"

time klog total --no-warn "${FILE}"
