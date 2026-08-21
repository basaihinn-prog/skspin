#!/bin/bash
# Shared deploy helpers used by CI workflows and local operators.
# Non-interactive by design: every function takes explicit args (no `select` prompts)
# so it is safe to run from GitHub Actions / Drone as well as by hand.
#
# Required env for remote operations:
#   DEPLOY_HOST       ssh alias or user@host of the target VPS
#   DEPLOY_SSH_OPTS   (optional) extra ssh options, e.g. "-i /path/key -o StrictHostKeyChecking=accept-new"
set -euo pipefail

SSH_OPTS="${DEPLOY_SSH_OPTS:--o StrictHostKeyChecking=accept-new}"

log() { echo "[deploy] $*" >&2; }
die() { echo "[deploy][ERROR] $*" >&2; exit 1; }

require_host() {
  [[ -n "${DEPLOY_HOST:-}" ]] || die "DEPLOY_HOST is not set (pass the target VPS ssh alias/user@host)"
}

# rsync_push <local_path> <remote_path> [extra rsync args...]
rsync_push() {
  local local_path="$1"; local remote_path="$2"; shift 2
  require_host
  log "rsync ${local_path} -> ${DEPLOY_HOST}:${remote_path}"
  rsync -avzP -e "ssh ${SSH_OPTS}" "$@" "${local_path}" "${DEPLOY_HOST}:${remote_path}"
}

remote_run() {
  require_host
  # shellcheck disable=SC2029
  ssh ${SSH_OPTS} "${DEPLOY_HOST}" "$@"
}

# deploy_binary <local_bin_path> <service_name> <remote_bin_dir default=/data/game/bin>
# Ships a single binary as "<service>.new", then does an atomic swap + restart via
# the existing `launch` process supervisor already installed on every VPS
# (see gameadmin/backend/banjia.md). Verifies the process is alive afterwards and
# rolls back automatically if it is not.
deploy_binary() {
  local local_bin="$1" service="$2" remote_dir="${3:-/data/game/bin}"
  require_host
  [[ -f "$local_bin" ]] || die "binary not found: $local_bin"

  rsync_push "$local_bin" "${remote_dir}/${service}.new"

  log "restarting ${service} on ${DEPLOY_HOST}"
  remote_run "set -e
    cd '${remote_dir}'
    pkill -x '${service}' || true
    sleep 1
    if [[ -f '${service}' ]]; then mv '${service}' '${service}.bak'; fi
    mv '${service}.new' '${service}'
    chmod +x '${service}'
    launch './${service}'
    sleep 2
    if pgrep -x '${service}' >/dev/null; then
      echo 'OK: ${service} is running'
    else
      echo 'FAIL: ${service} did not start, rolling back' >&2
      if [[ -f '${service}.bak' ]]; then
        mv '${service}.bak' '${service}'
        launch './${service}'
      fi
      exit 1
    fi"
}

# deploy_frontend_dist <local_dist_dir> <remote_dir>
deploy_frontend_dist() {
  local local_dist="$1" remote_dir="$2"
  require_host
  [[ -d "$local_dist" ]] || die "dist dir not found: $local_dist"
  rsync_push "${local_dist%/}/" "${remote_dir}" --delete
}

# deploy_static_assets <local_dir> <remote_dir> [--delete]
deploy_static_assets() {
  local local_dir="$1" remote_dir="$2"; shift 2 || true
  require_host
  rsync_push "${local_dir%/}/" "${remote_dir}" "$@"
}
