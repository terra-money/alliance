package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	custombankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	govtypes "github.com/terra-money/alliance/custom/gov/types"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
)

type Keeper struct {
	govkeeper.Keeper

	ak       alliancekeeper.Keeper
	custombk custombankkeeper.Keeper
	customsk govtypes.StakingKeeper
	bk       bankkeeper.Keeper
	sk       govtypes.StakingKeeper
	acck     accountkeeper.AccountKeeper
}

var _ govkeeper.Keeper = Keeper{}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	paramSpace types.ParamSubspace,
	ak accountkeeper.AccountKeeper,
	bk bankkeeper.BaseKeeper,
	sk stakingkeeper.Keeper,
	legacyRouter v1beta1.Router,
	router *baseapp.MsgServiceRouter,
	config types.Config,
) Keeper {
	keeper := Keeper{
		Keeper: govkeeper.NewKeeper(cdc, key, paramSpace, ak, bk, sk, legacyRouter, router, config),
		ak:     alliancekeeper.Keeper{},
		bk:     custombankkeeper.Keeper{},
		sk:     stakingkeeper.Keeper{},
		acck:   ak,
	}
	return keeper
}

// func (k *Keeper) RegisterKeepers(ak alliancekeeper.Keeper, bk custombankkeeper.Keeper, sk govtypes.StakingKeeper) {
// 	k.ak = ak
// 	k.bk = bk
// 	k.sk = sk
// }
