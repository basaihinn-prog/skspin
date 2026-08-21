#!/bin/bash
# Non-interactive sync: pass the target host as arg 1 or set DEPLOY_HOST.
# Kept non-interactive so this runs unattended from Drone/GitHub Actions.
set -euo pipefail

host="${1:-${DEPLOY_HOST:-}}"
valid_hosts="doudou-test doudou-prod dou-ph-prod dou-idr-prod"

if [[ -z "$host" ]]; then
  echo "ERROR: target host required. Usage: ./sync.sh <host> (one of: ${valid_hosts})" >&2
  exit 1
fi

case " ${valid_hosts} " in
  *" ${host} "*) ;;
  *) echo "ERROR: unknown host '${host}'. Expected one of: ${valid_hosts}" >&2; exit 1 ;;
esac

echo "syncing dist/ -> ${host}:/data/game-admin/"
rsync -avz -P --delete ./dist/ "root@${host}:/data/game-admin/"
