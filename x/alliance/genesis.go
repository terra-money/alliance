package alliance

import (
	"alliance/x/alliance/types"
	"time"
)

// ValidateGenesis
func ValidateGenesis(data *types.GenesisState) error {
	return nil
}

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:     24 * 60 * 60 * 1000_000_000,
			RewardClaimInterval: 5 * 60 * 1000_000_000,
			LastRewardClaimTime: time.Now(),
		},
		Assets: []types.AllianceAsset{},
	}
}
