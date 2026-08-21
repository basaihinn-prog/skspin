#!/bin/bash
# Non-interactive dist sync for the player-facing site (skyxspin variant).
# Usage: ./sync.sh <host> <remote_dir>
set -euo pipefail

host="${1:-${DEPLOY_HOST:-}}"
remote_dir="${2:-${DEPLOY_REMOTE_DIR:-}}"
valid_hosts="doudou-test doudou-prod dou-ph-prod dou-idr-prod"

if [[ -z "$host" || -z "$remote_dir" ]]; then
  echo "ERROR: usage: ./sync.sh <host> <remote_dir>  (hosts: ${valid_hosts})" >&2
  exit 1
fi

case " ${valid_hosts} " in
  *" ${host} "*) ;;
  *) echo "ERROR: unknown host '${host}'. Expected one of: ${valid_hosts}" >&2; exit 1 ;;
esac

echo "syncing dist/ -> ${host}:${remote_dir}"
rsync -avz -P --delete ./dist/ "root@${host}:${remote_dir}"
