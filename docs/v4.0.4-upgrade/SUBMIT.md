# Passage v4.0.4 software-upgrade proposal — submit guide

**Upgrade name:** `v4.0.4` (must match `app/upgrades/v4.0.4` `const Name`)
**Height:** `PENDING_HEIGHT` (pick ~7 days of voting + margin at ~6.05s blocks; recompute the day of submit)
**Deposit:** 100,000 PASG (`100000000000upasg`)

## What it changes
Binary-only fix: Amino JSON (Ledger, Keplr/Leap Amino) signing of CosmWasm
messages. Backports CosmWasm/wasmd#1844 + #1860 onto wasmd v0.50.0 via
`./wasmd-fork` (see `wasmd-fork/PASSAGE_PATCH.md`). No store, param or state
migration; `StoreUpgrades` is empty and the handler only runs `RunMigrations`.

## Why it needs a coordinated upgrade (not a rolling restart)
It changes which transactions pass signature verification. If some validators
ran v4.0.4 and others v4.0.3, an Amino-signed `MsgExecuteContract` in a block
would succeed on one set and fail (code 4) on the other, giving different tx
results / `LastResultsHash` and halting the chain. Every validator must switch
at the same height.

## Before submit
1. Tag and publish the v4.0.4 release; replace `PENDING_AMD64` / `PENDING_ARM64`
   in `proposal_v4.0.4.json` with the real sha256 checksums.
2. Fork-test on mainnet state: node logs `applying upgrade "v4.0.4"` and keeps
   producing blocks; all wasm codes/contracts preserved.
3. Smoke-test after the upgrade height on the fork: a Ledger (Amino) signed
   `MsgExecuteContract` with and without `funds` is accepted (code 0).
4. Set `height` and re-verify against the current chain height the day of submit.

## Submit (gov v1, on a node with the signing key)
```
passage tx gov submit-proposal docs/v4.0.4-upgrade/proposal_v4.0.4.json \
  --from <key> --chain-id passage-2 \
  --gas auto --gas-adjustment 1.5 --gas-prices 12.5upasg \
  --node <rpc> -y
```
Then deposit (if not met) + vote. Verify on-chain plan name reads exactly `v4.0.4`.
