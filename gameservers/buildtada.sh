#!/bin/bash
# Mutates servicejili source (jili_ -> tada_) before building, matching the
# behavior of buildtada.bat. Must run BEFORE buildjili.sh in buildall.sh.
set -euo pipefail
cd "$(dirname "$0")"

export GOOS=linux
export GOARCH=arm64
serviceDir=servicejili
outputDir=bin

mkdir -p "$outputDir"

sed -i 's/jili:/tada:/g' "${serviceDir}/jiliut/regRpc.go"

for dir in servicejili/*/; do
  service="$(basename "$dir")"
  echo "$dir"
  sed -i 's/jili_/tada_/g' "${serviceDir}/${service}/internal/const.go"
  out="${service//jili_/tada_}"
  go build -trimpath -ldflags "-s -w" -o "${outputDir}/${out}" "./${serviceDir}/${service}"
done

if [[ "${1:-}" != "skip" ]]; then
  ( cd "$outputDir" && 7z a bin * )
fi
