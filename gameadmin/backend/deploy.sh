#!/bin/bash
# Non-interactive build + deploy for a single gameadmin/backend service.
# Replaces the interactive `select host` prompt in sync-restart.sh so this can
# run unattended from CI (GitHub Actions / Drone) as well as by hand.
#
# Usage: DEPLOY_HOST=doudou-test ./deploy.sh <service> [AUTHOR]
#   or:  ./deploy.sh <service> <host> [AUTHOR]
set -euo pipefail
cd "$(dirname "$0")"

SERVICE="${1:?usage: deploy.sh <service> [host] [author]}"
HOST_ARG="${2:-${DEPLOY_HOST:-}}"
AUTHOR="${3:-${DEPLOY_AUTHOR:-$(git config user.name 2>/dev/null || echo ci)}}"

[[ -n "$HOST_ARG" ]] || { echo "ERROR: target host required (arg 2 or DEPLOY_HOST env)" >&2; exit 1; }
export DEPLOY_HOST="$HOST_ARG"

source ../../scripts/deploy/common.sh

log "building ${SERVICE} (author=${AUTHOR})"
sh build.sh "$SERVICE" "bin/${SERVICE}" "$AUTHOR"

deploy_binary "bin/${SERVICE}" "$SERVICE" "/data/game/bin"
