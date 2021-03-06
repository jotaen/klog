#!/bin/sh

set -e

echo 'Print help text'
klog --help

echo 'Create sample file'
FILE='time.klg'
echo '
2020-01-15
Did #something
  1h this
  13:00-14:00 that
' > "${FILE}"
cat "${FILE}"

echo 'Evaluate sample file'
klog total "${FILE}"
