#!/bin/bash
# Build gameservers engines and deploy binaries to /data/game/bin/ on the
# target host (same rsync+launch mechanism as gameadmin/backend).
#
# WARNING: with no service filter this rebuilds and restarts ALL ~400 game
# engines sequentially - each restart briefly stops that engine's process.
# Prefer passing explicit service names for targeted redeploys:
#
# Usage: ./deploy.sh <host> [service...]
#   ./deploy.sh doudou-test                 # build + deploy every engine (slow, use sparingly)
#   ./deploy.sh doudou-test service_pg_100  # build + deploy a single engine
set -euo pipefail
cd "$(dirname "$0")"

HOST_ARG="${1:?usage: deploy.sh <host> [service...]}"
shift
export DEPLOY_HOST="$HOST_ARG"

source ../scripts/deploy/common.sh

log "building gameservers (linux/arm64)"
sh buildall.sh

if [[ $# -gt 0 ]]; then
  targets=("$@")
else
  targets=($(cd bin && ls))
fi

for svc in "${targets[@]}"; do
  deploy_binary "bin/${svc}" "${svc}" "/data/game/bin"
done
