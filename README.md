# Passage3D

Passage is a metaverse platform on Cosmos: 3D worlds streamed through the browser, with NFT
collections and a marketplace for in-world assets. This repository is the Passage blockchain,
`passage`: a Cosmos SDK application with CosmWasm smart contracts, IBC and the Passage airdrop
module (`x/claim`).

## Network

| | |
|---|---|
| Chain ID | `passage-2` |
| Binary | `passage` |
| Denom | `upasg` (1 PASG = 10^6 upasg) |
| Address prefix | `pasg` |
| Minimum gas price | `12.5upasg` (nodes enforce this floor) |
| Mainnet version | v4.0.3; next upgrade `v4.0.4` ([proposal guide](docs/v4.0.4-upgrade/SUBMIT.md)) |

For running a validator or full node see [How to run a node](https://github.com/envadiv/mainnet#join-as-a-validator).
Set `minimum-gas-prices = "12.5upasg"` (or higher) in `app.toml`.

## Stack

| Component | Version |
|---|---|
| Cosmos SDK | v0.50.13 |
| CometBFT | v0.38.23 |
| IBC | ibc-go v8.0.0 |
| CosmWasm | wasmd v0.50.0 with the Amino JSON fixes of wasmd v0.51 ([`wasmd-fork/`](wasmd-fork/PASSAGE_PATCH.md)), wasmvm v1.5.7 |
| IAVL | v1.2.3 with the empty-store fixes ([`iavl-fork/`](iavl-fork/passage_empty_store_test.go)) |
| Go | 1.22.11 |

## Repository layout

| Path | Content |
|---|---|
| `cmd/passage` | Node and CLI |
| `app` | Application wiring, ante handlers, upgrade handlers (`app/upgrades/<version>`) |
| `x/claim` | Airdrop claim module (`MsgClaim`, claim records and queries) |
| `proto` | Protobuf definitions of the Passage modules |
| `wasmd-fork`, `iavl-fork` | Patched dependencies, wired in with `replace` in `go.mod` |
| `scripts/tests` | Docker-based upgrade, localnet and contract tests |
| `docs` | Upgrade proposals and submit guides |
| `.github/workflows` | CI and release pipeline |

## Build

Requirements: Go 1.22.11 and gcc (Ledger support). `LEDGER_ENABLED=false` builds without gcc.

```sh
make build            # ./build/passage
make install          # $GOPATH/bin/passage
make build-linux      # static linux binary (musl + libwasmvm), built in Docker
make docker-build     # runtime image passage:local (ports 26656, 26657, 1317, 9090)
```

## Development

| Target | What it does |
|--------|--------------|
| `make test` | Unit tests |
| `make test-forks` | Passage tests of the vendored `iavl-fork` and `app/params` |
| `make test-sim-*` | Application simulations (full, import/export, after-import, non-determinism) |
| `make test-upgrade` | In-place upgrade from the mainnet release (v4.0.3) to the working tree, two nodes, in Docker |
| `make test-contracts` | CosmWasm contracts across the same upgrade, in Docker (see below) |
| `make lint` | golangci-lint |
| `make proto-gen` / `make proto-lint` / `make proto-check-breaking` | Protobuf code generation, lint and breaking-change check (Docker) |
| `make localnet-start` / `localnet-stop` / `localnet-clean` | 4-validator localnet (`docker-compose.yml`, configuration in `./mytestnet`) |

### Contract tests

`make test-contracts` stores the `passage-contracts` artifacts (legacy contracts excluded) and the
contract code deployed on passage-2 on a two-node local chain:

1. It runs every scenario on v4.0.3: minting, marketplace sales with fee and royalty payouts,
   auctions, staking vault, royalty groups, follow/friend, and the cw3/cw4 multisig.
2. It leaves open orders, auctions and proposals, then upgrades the chain to the working tree.
3. It checks that the raw state of every contract is unchanged and completes the open operations.
4. It runs the scenarios again and compares the v4.0.3 and new outcomes.

| Variable | Use |
|---|---|
| `CONTRACTS_DIR` | Contract artifacts. Default: `../passage-workspace/passage-contracts/artifacts` |
| `MAINNET_DIR` | Directory with `code_<id>.wasm` files. Default: download the passage-2 code over gRPC |
| `QUICK=1` | Run the scenarios on the working tree only, without the upgrade |

Results are written to `build/contracts-test/`: `summary.txt`, `compare.txt` and the node logs.

### CI

Pull requests run these workflows:

- build and Docker image
- unit tests and the vendored fork tests
- the in-place upgrade test
- simulations
- golangci-lint
- CodeQL
- buf lint and breaking-change check, when protobuf files change

## Releases

Pushing a tag `vX.Y.Z` (or `vX.Y.Z-rcN`) runs `.github/workflows/release.yml`. GoReleaser builds
static `passage` binaries for linux amd64 and arm64 and publishes them with `sha256sum.txt`.
The checksums go into the upgrade proposal (`docs/<version>-upgrade/`).

## Contributing

- Use conventional commit titles (`feat`, `fix`, `chore`, `docs`, `test`, `ci`, ...).
- Put state-machine changes behind an upgrade handler in `app/upgrades/` and cover them with tests.
- Run `make test` and `make lint`. Also run `make test-upgrade` / `make test-contracts` when a
  change affects the upgrade path or CosmWasm.
- Record each change in `CHANGELOG.md`.

## License

[Apache License 2.0](LICENSE)
