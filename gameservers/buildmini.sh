#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"

export GOOS=linux
export GOARCH=arm64
outputDir=bin

mkdir -p "$outputDir"

for dir in serviceSpribeGame/*/; do
  service="$(basename "$dir")"
  echo "$dir"
  go build -trimpath -ldflags "-s -w" -o "${outputDir}/${service}" "./serviceSpribeGame/${service}"
done

if [[ "${1:-}" != "skip" ]]; then
  ( cd "$outputDir" && 7z a bin * )
fi
