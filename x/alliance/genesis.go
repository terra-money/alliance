package alliance

import (
	"alliance/x/alliance/types"
)

// ValidateGenesis
func ValidateGenesis(data *types.GenesisState) error {
	return nil
}

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{
			RewardDelayTime: 24 * 60 * 60 * 1000_000_000,
		},
		Assets: []types.AllianceAsset{},
	}
}
