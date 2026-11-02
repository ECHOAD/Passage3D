package v4_0_4

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/envadiv/Passage3D/app/upgrades"
	claim "github.com/envadiv/Passage3D/x/claim/keeper"
)

// Name is the on-chain upgrade name for Passage v4.0.4.
//
// v4.0.4 is a consensus-breaking release; every validator must switch binaries at
// the upgrade height. It contains:
//
//   - Amino JSON (SIGN_MODE_LEGACY_AMINO_JSON) signing of CosmWasm messages:
//     backport of CosmWasm/wasmd#1844 and #1860 onto wasmd v0.50.0 (see
//     wasmd-fork/PASSAGE_PATCH.md). Ledger / Keplr Amino signatures over
//     MsgExecuteContract & co. are accepted again.
//   - x/claim: MsgClaim carries the cosmos.msg.v1.signer annotation (claim txs
//     were rejected by every node since v4.0.0); the module account is created with
//     a unique account number; an expired airdrop is disabled instead of being
//     "ended" again on every block.
//   - x/gov: the legacy (v1beta1) proposal router is registered again (text and
//     x/params ParameterChangeProposal content).
//   - ante: wasm GasRegisterDecorator; the blocked community-pool multisig can no
//     longer bypass the restriction by adding a MsgFundCommunityPool to the tx.
//   - CosmWasm VM: wasmvm v1.5.0 -> v1.5.7 (same version as v2.6.0 / v3.0.0).
//
// State migration (this handler): duplicate account numbers are repaired, see
// FixDuplicateAccountNumbers. Module consensus versions are unchanged, so
// RunMigrations is a no-op.
//
// IMPORTANT: this string MUST exactly equal the plan.name in the gov proposal.
const Name = "v4.0.4"

// Upgrade: no new/renamed/deleted module stores, so StoreUpgrades is empty.
var Upgrade = upgrades.Upgrade{
	UpgradeName:          Name,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

// CreateUpgradeHandler returns the v4.0.4 upgrade handler.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ distribution.Keeper,
	_ bank.Keeper,
	ak auth.AccountKeeper,
	_ claim.Keeper,
	_ consensusparamkeeper.Keeper,
	_ paramskeeper.Keeper,
	_ *stakingkeeper.Keeper,
	_ govkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		logger := sdk.UnwrapSDKContext(ctx).Logger().With("upgrade", Name)

		changed, err := FixDuplicateAccountNumbers(ctx, ak)
		if err != nil {
			return nil, fmt.Errorf("fixing duplicate account numbers: %w", err)
		}
		logger.Info("repaired account numbers", "accounts", changed)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
