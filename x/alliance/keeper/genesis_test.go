package keeper_test

import (
	"alliance/x/alliance/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	app, ctx := createTestContext(t)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:     time.Duration(1000000),
			RewardClaimInterval: time.Duration(1000000),
			LastRewardClaimTime: time.Unix(0, 0).UTC(),
		},
		Assets: []types.AllianceAsset{
			{
				Denom:        "stake",
				RewardWeight: sdk.NewDec(1.0),
				TakeRate:     sdk.NewDec(0.0),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})

	delay := app.AllianceKeeper.RewardDelayTime(ctx)
	require.Equal(t, time.Duration(1000000), delay)

	interval := app.AllianceKeeper.RewardClaimInterval(ctx)
	require.Equal(t, time.Duration(1000000), interval)

	lastClaimTime := app.AllianceKeeper.LastRewardClaimTime(ctx)
	require.Equal(t, time.Unix(0, 0).UTC(), lastClaimTime)

	assets := app.AllianceKeeper.GetAllAssets(ctx)
	require.Equal(t, 1, len(assets))
	require.Equal(t, types.AllianceAsset{
		Denom:        "stake",
		RewardWeight: sdk.NewDec(1.0),
		TakeRate:     sdk.NewDec(0.0),
		TotalTokens:  sdk.ZeroInt(),
	}, assets[0])
}
