package alliance

import (
	"github.com/terra-money/alliance/x/alliance/types"
)

// ValidateGenesis
func ValidateGenesis(data *types.GenesisState) error {
	params := data.Params
	if params.TakeRateClaimInterval <= 0 {
		return types.ErrInvalidGenesisState.Wrap("reward_claim_interval has to be more than 0")
	}
	if len(data.Delegations) > 0 && len(data.Assets) == 0 {
		return types.ErrInvalidGenesisState.Wrap("cannot have delegations without alliance assets")
	}
	if len(data.Delegations) > 0 && len(data.ValidatorInfos) == 0 {
		return types.ErrInvalidGenesisState.Wrap("cannot have delegations without alliance validator infos")
	}
	if len(data.Redelegations) > 0 && len(data.Delegations) == 0 {
		return types.ErrInvalidGenesisState.Wrap("cannot have redelegations without delegations")
	}
	return nil
}

func DefaultGenesisState() *types.GenesisState {
	return &types.GenesisState{
		Params:                     types.DefaultParams(),
		Assets:                     []types.AllianceAsset{},
		ValidatorInfos:             []types.ValidatorInfoState{},
		RewardWeightChangeSnaphots: []types.RewardWeightChangeSnapshotState{},
		Delegations:                []types.Delegation{},
		Redelegations:              []types.RedelegationState{},
		Undelegations:              []types.UndelegationState{},
	}
}
