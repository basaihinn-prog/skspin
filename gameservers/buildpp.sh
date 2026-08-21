#!/bin/bash
# NOTE: buildpp.bat referenced an unset %out% variable (latent bug) - fixed
# here to use the service directory name as the output binary name.
set -euo pipefail
cd "$(dirname "$0")"

export GOOS=linux
export GOARCH=arm64
outputDir=bin

mkdir -p "$outputDir"

for dir in servicepp/*/; do
  service="$(basename "$dir")"
  echo "$dir"
  go build -trimpath -ldflags "-s -w" -o "${outputDir}/${service}" "./servicepp/${service}"
done

if [[ "${1:-}" != "skip" ]]; then
  ( cd "$outputDir" && 7z a bin * )
fi
