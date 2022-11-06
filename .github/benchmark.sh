#!/bin/bash

set -e

declare -a samples=(1 10 100 1000 10000 100000)
ITERATIONS=3

# Ensure binary is all set.
klog > /dev/null

# Run benchmark.
TIMEFORMAT=%R
for size in "${samples[@]}"; do
  printf "%8d:  " $size
  for _ in $(seq $ITERATIONS); do
    # Generate new test data.
    file="$(mktemp)"
    go run benchmark.go "${size}" > "${file}"

    runtime=$( { time klog total --no-warn "${file}" > /dev/null; } 2>&1 )
    printf "%ss  " $runtime
  done
  echo
done
