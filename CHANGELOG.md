<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState
given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

## v4.0.4 (consensus-breaking, on-chain upgrade `v4.0.4`)

Base: v4.0.3 (Cosmos SDK v0.50.13, CometBFT v0.38.23, wasmd v0.50.0, ibc-go v8.0.0).
Every validator must switch binaries at the upgrade height. See `docs/v4.0.4-upgrade/SUBMIT.md`.

### State machine

- CosmWasm: Amino JSON sign bytes fixed (backport of CosmWasm/wasmd#1844 and #1860 in `./wasmd-fork`);
  Ledger / Keplr / Leap Amino / CosmJS Amino-signed wasm transactions are accepted again.
- CosmWasm VM: `github.com/CosmWasm/wasmvm` v1.5.0 -> v1.5.7 (same as v2.6.0 / v3.0.0; security releases).
- x/gov: legacy (v1beta1) proposal router restored (text, x/params `ParameterChangeProposal`).
- x/claim: `MsgClaim` has the `cosmos.msg.v1.signer` and `amino.name` ("claim/Claim") annotations
  (claim txs were rejected since v4.0.0); unknown claim actions rejected; the module account is
  created with a unique account number; an expired airdrop is disabled instead of `EndAirdrop` on every block.
- Upgrade handler: repairs duplicate account numbers (claim module account vs account number 0 on passage-2).
- Ante: wasm `GasRegisterDecorator` (authz exec of gzipped `MsgStoreCode`); the blocked community-pool
  multisig can no longer bypass its restriction by adding a `MsgFundCommunityPool` to the tx;
  CheckTx minimum gas price floor (12.5upasg) restored.

### Node and CLI

- `export --for-zero-height` ported to SDK v0.50 (was a panic); `--modules-to-export` honoured.
- iavl-fork `hasVersion`: empty stores no longer report missing versions (rollback failed with
  "version does not exist" and left the node unable to restart; historical queries failed).
- Nodes apply the 12.5upasg minimum gas price floor (warning) when `minimum-gas-prices` is lower;
  new `app.toml` files default to 12.5upasg.
- REST: node service routes (`/cosmos/base/node/v1beta1/config`, `/status`) and swagger registered again.
- CLI: `SIGN_MODE_TEXTUAL`, `config` (confix), `prune`, `snapshots`, `genesis`, `q block(s)`,
  `q block-results`, `tx simulate`; `[wasm]` app.toml section and `--wasm.*` start flags with libwasmvm
  version check; no duplicated module commands; temp directories cleaned; CLI errors printed.
- `passage init` / `passage testnet` produce `upasg` genesis files; testnet uses 5% min commission.

### Build, release, tests

- Dockerfile: Go 1.22.11, static musl binary with checksum-verified libwasmvm v1.5.7, Ledger support.
- Makefile: version from git (fallback to commit), fixed build-tags stamp, proto via proto-builder 0.14.0,
  real simulation targets, `localnet-*` on Docker, `test-upgrade`.
- Release workflow: `GOTOOLCHAIN=go1.22.11` (v4 releases failed with goreleaser-cross Go 1.22.7);
  CI workflows updated to Go from `go.mod`, branch `main`.
- Test suite ported to SDK v0.50 (it no longer compiled since v4.0.0); new tests for every fix,
  application simulations (`app/sim_test.go`), in-place upgrade test from v4.0.3
  (`scripts/tests/upgrade_test.sh`) and a 4-validator localnet scenario test (`scripts/tests/localnet_test.sh`).

## [v1.1.0](https://github.com/envadiv/Passage3D/releases/tag/v1.1.0) - 2022-10-17

- Update cosmos-sdk to v0.45.9 and tendermint v0.34.21 which includes the patch for IBC issue
- fix lint by @anilCSE in https://github.com/envadiv/Passage3D/pull/79

**Full Changelog**: https://github.com/envadiv/Passage3D/compare/v1.0.0...v1.1.0

## [v1.0.0](https://github.com/envadiv/Passage3D/releases/tag/v1.0.0) - 2022-08-08

Initial Release!
