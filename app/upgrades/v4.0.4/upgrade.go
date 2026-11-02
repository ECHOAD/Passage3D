package v4_0_4

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
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
// v4.0.4 fixes Amino JSON (SIGN_MODE_LEGACY_AMINO_JSON) signing of CosmWasm
// messages by backporting CosmWasm/wasmd#1844 onto wasmd v0.50.0 (see
// wasmd-fork/PASSAGE_PATCH.md). The change is binary-only, but it changes which
// transactions pass signature verification, so every validator must switch
// binaries at the same height: a node still on v4.0.3 would reject an
// Amino-signed MsgExecuteContract that a v4.0.4 node accepts, producing
// different tx results (LastResultsHash) and halting consensus.
//
// IMPORTANT: this string MUST exactly equal the plan.name in the gov proposal.
const Name = "v4.0.4"

// Upgrade: no new/renamed/deleted module stores, so StoreUpgrades is empty.
var Upgrade = upgrades.Upgrade{
	UpgradeName:          Name,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

// CreateUpgradeHandler matches the shared upgrades.Upgrade signature. All
// keepers are unused: v4.0.4 makes no state or param changes.
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ distribution.Keeper,
	_ bank.Keeper,
	_ auth.AccountKeeper,
	_ claim.Keeper,
	_ consensusparamkeeper.Keeper,
	_ paramskeeper.Keeper,
	_ *stakingkeeper.Keeper,
	_ govkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
