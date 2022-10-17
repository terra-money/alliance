package keeper

import (
	"testing"

	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func AllianceKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	legacyAmino := codec.NewLegacyAmino()

	ps := typesparams.NewSubspace(
		cdc,
		legacyAmino,
		storeKey,
		memStoreKey,
		"AllianceParams",
	)

	ctrl := gomock.NewController(t)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		ps,
		NewMockAccountKeeper(ctrl),
		NewMockBankKeeper(ctrl),
		NewMockStakingKeeper(ctrl),
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

func CreateNewAllianceAsset(Keeper *keeper.Keeper, ctx sdk.Context, n int64) types.AllianceAsset {
	return types.AllianceAsset{
		Denom:        "uluna",
		RewardWeight: sdk.NewDec(1),
		TakeRate:     sdk.NewDec(1),
		TotalTokens:  math.NewInt(10 * n),
		TotalShares:  sdk.NewDec(10).Mul(sdk.NewDec(n)),
	}
}

func CreateNewDelegation(Keeper *keeper.Keeper, ctx sdk.Context, n int64) []types.Delegation {
	items := make([]types.Delegation, n)
	for i := range items {
		items[i] = types.Delegation{
			DelegatorAddress: "alliance1reanw526y9nun5mtajwkf6p8j3g8c6x276hvfl",
			ValidatorAddress: "alliancevaloper1reanw526y9nun5mtajwkf6p8j3g8c6x2gyesnr",
			Denom:            "uluna",
			Shares:           sdk.NewDec(10),
		}
	}
	return items
}
