#!/bin/bash

set -eou pipefail

SAM_DIR="$(dirname "$0")"
FUNC_DIR_RELATIVE="../function"

cd "${SAM_DIR}"

SANITY="${FUNC_DIR_RELATIVE}/main.go"
if [ ! -f "${SANITY}" ] ; then
  echo "ERROR: File not found: ${SANITY}" 1>&2
  exit 1
fi

cd "${FUNC_DIR_RELATIVE}"
GOOS=linux GOARCH=amd64 go build -o main
cd -

mv "${FUNC_DIR_RELATIVE}/main" ./
sam deploy \
  --no-confirm-changeset \
  --parameter-overrides \
  "ParameterKey=Foo,ParameterValue=${FOO}" \
  $@
rm -f ./main



