#!/bin/bash

# Performs a test of the built binary. That includes:
# - Doing a brief smoke test to check that the binary can be invoked
# - Checking that all build-time information got compiled in correctly

set -e

THIS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

echo 'Print help text...'
klog --help 1>/dev/null

echo 'Create sample file...'
FILE='time.klg'
echo '
2020-01-15
Did #something
  1h this
  13:00-14:00 that
' > "${FILE}"

echo 'Evaluate sample file...'
klog total "${FILE}" 1>/dev/null
[[ "$(klog total --no-style "${FILE}" | head -n 1)" == "Total: 2h" ]] || exit 1

echo 'Check version...'
ACTUAL_VERSION="$(klog version --no-check --quiet)"
EXPECTED_VERSION="$1"
[[ "${ACTUAL_VERSION}" == "${EXPECTED_VERSION}" ]] || exit 1

echo 'Check build hash...'
ACTUAL_BUILD_HASH="$(klog version --no-check | grep -oE '\[[abcdef0123456789]{7}]')"
EXPECTED_BUILD_HASH="$2"
[[ "${ACTUAL_BUILD_HASH}" == "[${EXPECTED_BUILD_HASH::7}]" ]] || exit 1

echo 'Check embedded license file...'
ACTUAL_LICENSE="$(klog info --license)"
EXPECTED_LICENSE="$(cat "${THIS_DIR}/../LICENSE.txt")"
[[ "${ACTUAL_LICENSE}" == "${EXPECTED_LICENSE}" ]] || exit 1
