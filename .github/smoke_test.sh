#!/bin/bash

# Performs a test of the built binary. That includes:
# - Doing a brief smoke test to check that the binary can be invoked
# - Checking that all build-time information got compiled in correctly

set -e

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

echo 'Check version...'
ACTUAL_VERSION="$(klog version --no-check --quiet)"
[[ "${ACTUAL_VERSION}" == "${EXPECTED_VERSION}" ]] || exit 1

echo 'Check build hash...'
ACTUAL_BUILD_HASH="$(klog version --no-check | grep -oE '\[[abcdef0123456789]{7}]')"
[[ "${ACTUAL_BUILD_HASH}" == "[${EXPECTED_BUILD_HASH::7}]" ]] || exit 1

echo 'Check embedded spec file...'
ACTUAL_SPEC="$(klog info --spec)"
[[ "${ACTUAL_SPEC}" == "$(cat "${EXPECTED_SPEC_PATH}")" ]] || exit 1

echo 'Check embedded license file...'
ACTUAL_LICENSE="$(klog info --license)"
[[ "${ACTUAL_LICENSE}" == "$(cat "${EXPECTED_LICENSE_PATH}")" ]] || exit 1
