#!/bin/bash

set -euxo pipefail

curl -sSfL https://raw.githubusercontent.com/getsentry/sentry/master/src/sentry/models/project.py \
  | awk '/GETTING_STARTED_DOCS_PLATFORMS = \[/ { platforms = 1; next } /\]/ { platforms = 0 } platforms { gsub(/[ ",]/, ""); print }' \
  > platforms.txt