package alliance

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateGenesis
func ValidateGenesis(data *types.GenesisState) error {
	return nil
}

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{
			RewardDelayTime: 24 * 60 * 60 * 1000_000_000,
			GlobalIndex:     sdk.NewDec(0),
		},
		Assets: []types.AllianceAsset{},
	}
}
