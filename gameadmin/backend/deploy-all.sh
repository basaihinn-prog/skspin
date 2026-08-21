#!/bin/bash
# Non-interactive replacement for banjia.sh / banjia-test.sh.
# Builds every service in buildall.sh and ships bin/ + *.sh to the target host.
# The old scripts hardcoded the target host via commented-out lines; here the
# host is always an explicit argument so CI can target any environment safely.
#
# Usage: ./deploy-all.sh <host>   e.g. ./deploy-all.sh doudou-test
set -euo pipefail
cd "$(dirname "$0")"

HOST_ARG="${1:?usage: deploy-all.sh <host> (doudou-test|doudou-prod|dou-ph-prod|dou-idr-prod)}"
export DEPLOY_HOST="$HOST_ARG"

source ../../scripts/deploy/common.sh

log "building all admin backend services"
sh buildall.sh

log "shipping bin/ and *.sh to ${DEPLOY_HOST}:/data/game/"
rsync_push "bin" "/data/game/" --exclude-from=banjia-exclude.txt
for f in *.sh; do
  rsync_push "$f" "/data/game/${f}"
done
