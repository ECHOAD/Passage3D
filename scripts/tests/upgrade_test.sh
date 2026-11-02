#!/usr/bin/env bash
# In-place upgrade regression test, run in Docker:
#   OLD binary = OLD_REF (default eea1f3e = v4.0.3, the passage-2 mainnet release)
#   NEW binary = the current working tree
# A validator and a full node run the old binary, state is created (bank, staking,
# authz, feegrant, group, CosmWasm, claim), a software-upgrade proposal for $UPG is
# passed, both nodes switch to the new binary and every v4.0.4 fix is checked.
#
# Usage: scripts/tests/upgrade_test.sh            (or: make test-upgrade)
# Results: build/upgrade-test/{summary.txt,upgrade.log,*.log}
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/../.." && pwd)
OUT=${OUT:-$ROOT/build/upgrade-test}
mkdir -p "$OUT"

docker run --rm \
  -e OLD_REF="${OLD_REF:-eea1f3e}" \
  -e UPG="${UPG:-v4.0.4}" \
  -e GOFLAGS=-buildvcs=false \
  -v "$ROOT":/src:ro \
  -v "$OUT":/out \
  ${GOMODCACHE_VOLUME:+-v "$GOMODCACHE_VOLUME":/go/pkg/mod} \
  golang:1.22.11-bookworm bash /src/scripts/tests/upgrade_inner.sh

grep -q '^FAIL' "$OUT/summary.txt" && { echo "upgrade test FAILED"; exit 1; }
echo "upgrade test passed"
