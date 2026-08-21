#!/bin/bash
# Build every gamecomponents service and deploy each binary to /data/game/bin/
# on the target host, using the same rsync+launch restart mechanism as
# gameadmin/backend (these binaries run alongside the admin services in the
# same /data/game/bin directory - see header.sh service list).
#
# Usage: ./deploy.sh <host> [service...]
#   ./deploy.sh doudou-test                 # build + deploy everything
#   ./deploy.sh doudou-test service_gateway # build + deploy a single service
set -euo pipefail
cd "$(dirname "$0")"

HOST_ARG="${1:?usage: deploy.sh <host> [service...]}"
shift
export DEPLOY_HOST="$HOST_ARG"

source ../scripts/deploy/common.sh

log "building gamecomponents (linux/arm64)"
sh build.sh skip

if [[ $# -gt 0 ]]; then
  targets=("$@")
else
  targets=($(cd bin && ls))
fi

for svc in "${targets[@]}"; do
  deploy_binary "bin/${svc}" "${svc}" "/data/game/bin"
done
