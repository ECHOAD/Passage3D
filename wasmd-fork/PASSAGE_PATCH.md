# wasmd v0.50.0 + Amino JSON fixes (Passage v4.0.4)

This directory is the `github.com/CosmWasm/wasmd@v0.50.0` module source, wired
in through `replace github.com/CosmWasm/wasmd => ./wasmd-fork` in the root
`go.mod` (same approach as `iavl-fork/`).

## What is patched

Two upstream fixes, both released in wasmd **v0.51.0**, backported without any
other v0.51 change (no wasmvm 2.0, no SDK bump):

| Upstream PR | Commit | Annotation | Effect on Amino JSON sign bytes |
|---|---|---|---|
| CosmWasm/wasmd#1844 "Fix Amino JSON representation" | `60375ab` | `(amino.encoding) = "inline_json"` on every `RawContractMessage` field | `msg` is rendered as inline (key-sorted) JSON instead of base64 |
| CosmWasm/wasmd#1860 "Add legacy_coins amino encoding annotation" | `5ec0c5e` | `(amino.encoding) = "legacy_coins"` on `funds` / `amount` coin lists | empty `funds` is rendered as `[]` instead of `null` |

Files replaced verbatim from upstream commit `5ec0c5e` (which contains both PRs):

- `proto/cosmwasm/wasm/v1/{tx,authz,proposal_legacy}.proto`
- `x/wasm/types/{tx,authz,proposal_legacy}.pb.go`

and from `60375ab` (types.proto is not touched by #1860):

- `proto/cosmwasm/wasm/v1/types.proto`, `x/wasm/types/types.pb.go`

These four files are byte-identical between `v0.50.0` and the parent of
`60375ab`, and no other commit touched them between `60375ab` and `5ec0c5e`, so
the result is exactly v0.50.0 + #1844 + #1860. In the `.pb.go` files only the
gzipped `FileDescriptorProto` bytes change: no Go type, field or protobuf wire
encoding changes. `query.proto` (also touched by #1844) is intentionally left
alone: queries are never signed, and its v0.51 version adds a new RPC.

## Why

Cosmos SDK v0.50 builds `SIGN_MODE_LEGACY_AMINO_JSON` sign bytes from protobuf
annotations (`cosmossdk.io/x/tx/signing/aminojson`) instead of the legacy amino
codec. Without the annotations above, a v0.50 node rebuilds
`MsgExecuteContract` as

    "funds":null, "msg":"eyJyZW1vdmVfYXNrIjp7InRva2VuX2lkIjoiMSJ9fQ=="

while Ledger, Keplr/Leap Amino signers and CosmJS/InterchainJS sign

    "funds":[],   "msg":{"remove_ask":{"token_id":"1"}}

so every Amino-signed CosmWasm transaction fails with
`signature verification failed ... unauthorized` (code 4). Direct
(`SIGN_MODE_DIRECT`) signing is unaffected.

`app/params/amino_contract_signing_test.go` pins the expected client bytes and
verifies a real secp256k1 signature over them; it fails on stock v0.50.0 and
passes with this fork.

## Removing the fork

Delete this directory and the `replace` line once Passage moves to wasmd
>= v0.51.0.
