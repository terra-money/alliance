package simulation

import (
	"alliance/x/alliance/types"
	"fmt"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"math/rand"
)

func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.RewardDelayTime),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genRewardDelayTime(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.RewardClaimInterval),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genRewardClaimInterval(r))
			},
		),
	}
}
