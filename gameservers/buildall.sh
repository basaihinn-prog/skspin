#!/bin/bash
# Linux equivalent of buildall.bat. Order matters: buildtada.sh then
# buildjili.sh must run in this sequence so servicejili/ source ends up back
# in the "jili" state that is committed to git (see comments in each script).
# Always restores the working tree afterwards so CI runs stay reproducible.
set -euo pipefail
cd "$(dirname "$0")"

rm -rf bin
sh build.sh skip
sh buildtada.sh skip
sh buildjili.sh skip
sh buildpp.sh skip
sh buildjdb.sh skip
sh buildhacksaw.sh skip
sh buildmini.sh

# buildtada.sh/buildjili.sh mutate tracked source files as a side effect of
# building (see comments there) - restore them so the checkout is clean.
git checkout -- servicejili/ 2>/dev/null || true
