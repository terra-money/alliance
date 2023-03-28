package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cometbft/cometbft/libs/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/terra-money/alliance/x/alliance/types"
)

func genRewardDelayTime(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

func genTakeRateClaimInterval(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 1, 60*60)) * time.Second
}

func genNumOfAllianceAssets(r *rand.Rand) int {
	return simulation.RandIntBetween(r, 0, 50)
}

func RandomizedGenesisState(simState *module.SimulationState) {
	var (
		rewardDelayTime     time.Duration
		rewardClaimInterval time.Duration
		numOfAllianceAssets int
	)

	r := simState.Rand
	rewardDelayTime = genRewardDelayTime(r)
	rewardClaimInterval = genTakeRateClaimInterval(r)
	numOfAllianceAssets = genNumOfAllianceAssets(r)

	var allianceAssets []types.AllianceAsset
	for i := 0; i < numOfAllianceAssets; i++ {
		rewardRate := simulation.RandomDecAmount(r, sdk.NewDec(5))
		takeRate := simulation.RandomDecAmount(r, sdk.MustNewDecFromStr("0.0005"))
		startTime := time.Now().Add(time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second)
		allianceAssets = append(allianceAssets, types.NewAllianceAsset(fmt.Sprintf("ASSET%d", i), rewardRate, sdk.NewDec(0), sdk.NewDec(15), takeRate, startTime))
	}

	allianceGenesis := types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       rewardDelayTime,
			TakeRateClaimInterval: rewardClaimInterval,
			LastTakeRateClaimTime: simState.GenTimestamp,
		},
		Assets: allianceAssets,
	}

	bz, err := json.MarshalIndent(&allianceGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated alliance parameters:\n%s\n", bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&allianceGenesis)
}
