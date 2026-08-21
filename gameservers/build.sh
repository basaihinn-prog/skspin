#!/bin/bash
# Linux build script equivalent to build.bat (for CI runners / Linux devs).
# Usage: ./build.sh [skip]   ("skip" = don't 7z-archive bin/, used by buildall.sh)
set -euo pipefail
cd "$(dirname "$0")"

export GOOS=linux
export GOARCH=arm64
outputDir=bin

mkdir -p "$outputDir"

for dir in service/*/; do
  service="$(basename "$dir")"
  echo "$dir"
  go build -trimpath -ldflags "-s -w" -o "${outputDir}/${service}" "./service/${service}"
done

if [[ "${1:-}" != "skip" ]]; then
  ( cd "$outputDir" && 7z a bin * )
fi
