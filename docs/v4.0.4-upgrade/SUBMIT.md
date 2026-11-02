# Passage v4.0.4 software-upgrade proposal — submit guide

**Upgrade name:** `v4.0.4` (must match `app/upgrades/v4.0.4` `const Name`)
**Height:** `PENDING_HEIGHT` (pick ~7 days of voting + margin at ~6.05s blocks; recompute the day of submit)
**Deposit:** 100,000 PASG (`100000000000upasg`, the passage-2 `min_deposit`)
**Message:** `/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade` (gov v1), like proposals #23 (v4.0.0)
and #24 (v4.0.2). Do **not** wrap it in `MsgExecLegacyContent`: v4.0.0–v4.0.3 have no legacy
gov router, so such a proposal fails on submit.

## What it changes (consensus-breaking)

| Area | Change |
|---|---|
| CosmWasm Amino signing | Backport of CosmWasm/wasmd#1844 + #1860 onto wasmd v0.50.0 (`./wasmd-fork`, see `wasmd-fork/PASSAGE_PATCH.md`). Ledger / Keplr / Leap Amino / CosmJS Amino signatures over wasm messages are accepted again. |
| CosmWasm VM | wasmvm v1.5.0 → v1.5.7 (the version used up to v3.0.0; v1.5.1–v1.5.7 are security releases). |
| x/gov | Legacy proposal router restored (`TextProposal`, x/params `ParameterChangeProposal`). |
| x/claim | `MsgClaim` signer annotation (claim txs were rejected since v4.0.0); expired airdrop disabled once instead of `EndAirdrop` on every block; unknown claim actions rejected; module account created with a unique account number. |
| State migration | Upgrade handler repairs duplicate account numbers (the claim module account shares number 0 with `pasg1lel0s624jr9zsz4ml6yv9e5r4uzukfs7hwh22w`); `RunMigrations` is a no-op (no consensus-version bumps). |
| Ante | wasm `GasRegisterDecorator`; blocked multisig can no longer bypass the restriction with a leading `MsgFundCommunityPool`; CheckTx minimum gas price floor 12.5upasg. |

Non-consensus fixes shipped in the same binary: `export --for-zero-height`, `rollback` and
historical queries on empty stores (iavl-fork `hasVersion`), node service REST routes, swagger,
`SIGN_MODE_TEXTUAL` in the CLI, `config`/`prune`/`snapshots`/`genesis` commands, wasm config
template and `--wasm.*` start flags, node refuses to start with `minimum-gas-prices` below 12.5upasg.

**Operators:** make sure `app.toml` has `minimum-gas-prices >= 12.5upasg` before starting v4.0.4.

## Why it needs a coordinated upgrade (not a rolling restart)
It changes which transactions pass signature verification and ante checks, and the
upgrade handler writes state. Mixed v4.0.3 / v4.0.4 validators would compute different
results for the same block and halt the chain. Every validator must switch at the same height.

## Before submit
1. `make test`, the `test-sim-*` targets and `make test-upgrade` green on the release commit.
2. Tag and publish the v4.0.4 release; fill `PENDING_REPO`, `PENDING_AMD64`, `PENDING_ARM64`
   in `proposal_v4.0.4.json` with the real repository and sha256 checksums of the release assets.
3. Fork-test on a copy of mainnet state (recommended): node logs `applying upgrade "v4.0.4"`,
   keeps producing blocks, the claim module account gets a new account number, all wasm
   codes/contracts preserved.
4. Set `height` and re-verify against the current chain height the day of submit.

## Submit (gov v1, on a node with the signing key)
```
passage tx gov submit-proposal docs/v4.0.4-upgrade/proposal_v4.0.4.json \
  --from <key> --chain-id passage-2 \
  --gas auto --gas-adjustment 1.5 --gas-prices 12.5upasg \
  --node <rpc> -y
```
Then deposit (if not met) + vote. Verify on-chain plan name reads exactly `v4.0.4`:
`passage q upgrade plan --node <rpc>`.
