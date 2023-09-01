#!/bin/bash

set -eou pipefail

SAM_DIR="$(dirname "$0")"
FUNC_DIR_RELATIVE="../function"
APP="bootstrap"

cd "${SAM_DIR}"

SANITY="${FUNC_DIR_RELATIVE}/main.go"
if [ ! -f "${SANITY}" ] ; then
  echo "ERROR: File not found: ${SANITY}" 1>&2
  exit 1
fi

cd "${FUNC_DIR_RELATIVE}"
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o "${APP}"
zip "${APP}.zip" "${APP}"
rm -f "${APP}"
cd -

mv "${FUNC_DIR_RELATIVE}/${APP}.zip" ./
sam deploy --tags "project=cloudwatch_to_google_chat" $@
rm -f "${APP}.zip"



