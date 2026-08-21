#!/bin/bash
# Mutates servicejili source (tada_ -> jili_) before building, matching the
# behavior of buildjili.bat. Must run AFTER buildtada.sh in buildall.sh so the
# working tree ends up in the "jili" state that is checked into git.
set -euo pipefail
cd "$(dirname "$0")"

export GOOS=linux
export GOARCH=arm64
serviceDir=servicejili
outputDir=bin

mkdir -p "$outputDir"

sed -i 's/tada:/jili:/g' "${serviceDir}/jiliut/regRpc.go"

for dir in servicejili/*/; do
  service="$(basename "$dir")"
  echo "$dir"
  sed -i 's/tada_/jili_/g' "${serviceDir}/${service}/internal/const.go"
  out="${service//tada_/jili_}"
  go build -trimpath -ldflags "-s -w" -o "${outputDir}/${out}" "./${serviceDir}/${service}"
done

if [[ "${1:-}" != "skip" ]]; then
  ( cd "$outputDir" && 7z a bin * )
fi
