# Build & Deploy Automation

This repo previously had **no real CI/CD**: builds were done by hand via Windows
`.bat` files + 7-Zip, and deploys were interactive bash scripts (`select host`
prompts) run manually over SSH. This document describes the automated
replacement.

## What changed

| Area | Before | Now |
|---|---|---|
| gameadmin/backend | manual `sh buildall.sh` + `banjia.sh`/`sync-restart.sh` (interactive) | `deploy-all.sh <host>` / `deploy.sh <service> <host>` (non-interactive) + GitHub Actions |
| gameadmin/web | Drone CI, but `sync.sh` prompted interactively (would hang under Drone) | `sync.sh <host>` non-interactive; Drone now auto-deploys to `doudou-test`; GitHub Actions added for prod |
| gamecomponents | Windows `build.bat` only, no deploy script | `build.sh` (Linux/CI) + `deploy.sh <host> [service...]` |
| gameservers | Windows `build*.bat` only, no deploy script, source-mutating jili/tada build order | `build*.sh` equivalents (git-tracked mutation is auto-reverted) + `deploy.sh <host> [service...]` |
| web / web_skyxspin | Windows `build.bat` + manual copy, no sync script | `sync.sh <host> <remote_dir>` + GitHub Actions |

All deploys are now driven by [scripts/deploy/common.sh](scripts/deploy/common.sh),
which:
- ships a binary as `name.new` via rsync,
- does an atomic `pkill` → `mv` → `launch` swap on the remote (same supervisor
  already in use, see [gameadmin/backend/banjia.md](gameadmin/backend/banjia.md)),
- verifies the process is alive with `pgrep`, and **automatically rolls back**
  to the previous binary if the new one fails to start.

## CI/CD pipeline

GitHub Actions workflows live in `.github/workflows/`:

- `deploy-admin-backend.yml`
- `deploy-admin-frontend.yml`
- `deploy-gamecomponents.yml`
- `deploy-gameservers.yml` — **push only compiles** (400+ engines); actual
  deploys always require `workflow_dispatch` with an explicit `services` list,
  so a single push never restarts every running game engine at once.
- `deploy-web.yml`
- `deploy-web-skyxspin.yml`

Each workflow:
1. **Push to `main`** → builds and auto-deploys to the `staging` GitHub
   Environment (intended to point at `doudou-test`).
2. **Manual `workflow_dispatch`** → choose `doudou-prod` / `dou-ph-prod` /
   `dou-idr-prod` as the target GitHub Environment.

`gameadmin/web/.drone.yml` still runs the same build and now calls
`sync.sh doudou-test` automatically (staging only) — kept alongside GitHub
Actions per request, since it already existed.

## One-time setup required in GitHub

### 1. SSH deploy key
Generate a dedicated key pair (do not reuse a personal key):
```sh
ssh-keygen -t ed25519 -f deploy_key -C "ci-deploy" -N ""
```
- Add the **public** key to `~/.ssh/authorized_keys` on every target VPS
  (`doudou-test`, `doudou-prod`, `dou-ph-prod`, `dou-idr-prod`) for the user
  that owns `/data/game` (root, based on existing scripts).
- Add the **private** key contents as repository secret `DEPLOY_SSH_KEY`
  (Settings → Secrets and variables → Actions → New repository secret; or set
  it per-Environment if different hosts should use different keys).

### 2. GitHub Environments
Create these under Settings → Environments:

| Environment | Suggested protection | Variable `DEPLOY_HOST` | Variable `DEPLOY_REMOTE_DIR` (web/web_skyxspin only) |
|---|---|---|---|
| `staging` | none (auto-deploy) | `doudou-test` ssh user@ip | e.g. `/data/h5games/web` |
| `doudou-prod` | **required reviewers** | `47.236.107.220` (or ssh alias) | e.g. `/data/dl/web` |
| `dou-ph-prod` | **required reviewers** | ph prod ssh user@ip | e.g. `/data/dl/web` |
| `dou-idr-prod` | **required reviewers** | idr prod ssh user@ip | e.g. `/data/h5games/web` |

`DEPLOY_HOST` must be something the GitHub-hosted runner can resolve directly
(`user@ip` or a public DNS name) — the `doudou-test` / `dou-ph-prod` style SSH
config aliases used interactively by operators only exist in their local
`~/.ssh/config` and won't resolve on a runner.

The exact `DEPLOY_REMOTE_DIR` for `web`/`web_skyxspin` was **not established**
anywhere in this repo (unlike `gameadmin/web`, which is confirmed to deploy to
`/data/game-admin/`) — confirm the real path with whoever manages the Caddy
config before enabling those two workflows for real traffic.

## Known limitations / follow-ups

- `gameservers/service_fish/Dockerfile` bakes a git password into a build arg —
  should move to BuildKit secrets or a deploy token before containerizing it.
- No automated health checks beyond `pgrep` — consider an HTTP health check
  where services expose one.
- `gameservers` deploys always rebuild everything (`buildall.sh`) even when
  targeting a single engine, to keep the jili/tada source-mutation ordering
  correct; this is slower but avoids a class of build bugs.
