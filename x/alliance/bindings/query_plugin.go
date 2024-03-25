package bindings

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/bindings/types"
	"github.com/terra-money/alliance/x/alliance/keeper"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

type QueryPlugin struct {
	allianceKeeper keeper.Keeper
}

func NewAllianceQueryPlugin(keeper keeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		allianceKeeper: keeper,
	}
}

func CustomQuerier(q *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) (result []byte, err error) {
	return func(ctx sdk.Context, request json.RawMessage) (result []byte, err error) {
		var AllianceRequest types.AllianceQuery
		err = json.Unmarshal(request, &AllianceRequest)
		if err != nil {
			return nil, err
		}
		if AllianceRequest.Alliance != nil {
			return q.GetAlliance(ctx, AllianceRequest.Alliance.Denom)
		}
		if AllianceRequest.Delegation != nil {
			denom := AllianceRequest.Delegation.Denom
			delegator := AllianceRequest.Delegation.Delegator
			validator := AllianceRequest.Delegation.Validator

			return q.GetDelegation(ctx, denom, delegator, validator)
		}
		if AllianceRequest.DelegationRewards != nil {
			denom := AllianceRequest.DelegationRewards.Denom
			delegator := AllianceRequest.DelegationRewards.Delegator
			validator := AllianceRequest.DelegationRewards.Validator

			return q.GetDelegationRewards(ctx,
				denom,
				delegator,
				validator,
			)
		}
		return nil, fmt.Errorf("unknown query")
	}
}

func (q *QueryPlugin) GetAlliance(ctx sdk.Context, denom string) (res []byte, err error) {
	asset, found := q.allianceKeeper.GetAssetByDenom(ctx, denom)
	if !found {
		return nil, alliancetypes.ErrUnknownAsset
	}
	res, err = json.Marshal(types.AllianceResponse{
		Denom:                asset.Denom,
		RewardWeight:         asset.RewardWeight.String(),
		TakeRate:             asset.TakeRate.String(),
		TotalTokens:          asset.TotalTokens.String(),
		TotalValidatorShares: asset.TotalValidatorShares.String(),
		RewardStartTime:      uint64(asset.RewardStartTime.Nanosecond()),
		RewardChangeRate:     asset.RewardChangeRate.String(),
		LastRewardChangeTime: uint64(asset.LastRewardChangeTime.Nanosecond()),
		RewardWeightRange: types.RewardWeightRange{
			Min: asset.RewardWeightRange.Min.String(),
			Max: asset.RewardWeightRange.Max.String(),
		},
		IsInitialized: asset.IsInitialized,
	})
	return
}

func (q *QueryPlugin) GetDelegation(ctx sdk.Context, denom string, delegator string, validator string) (res []byte, err error) {
	delegatorAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return
	}
	delegation, found := q.allianceKeeper.GetDelegation(ctx, delegatorAddr, validatorAddr, denom)
	if !found {
		return nil, alliancetypes.ErrDelegationNotFound
	}
	asset, found := q.allianceKeeper.GetAssetByDenom(ctx, denom)
	if !found {
		return nil, alliancetypes.ErrUnknownAsset
	}

	allianceValidator, err := q.allianceKeeper.GetAllianceValidator(ctx, validatorAddr)
	if err != nil {
		return nil, err
	}
	balance := alliancetypes.GetDelegationTokens(delegation, allianceValidator, asset)
	res, err = json.Marshal(types.DelegationResponse{
		Delegator: delegation.DelegatorAddress,
		Validator: delegation.ValidatorAddress,
		Denom:     delegation.Denom,
		Amount: types.Coin{
			Denom:  balance.Denom,
			Amount: balance.Amount.String(),
		},
	})
	return res, err
}

func (q *QueryPlugin) GetDelegationRewards(ctx sdk.Context,
	denom string,
	delegator string,
	validator string,
) (res []byte, err error) {
	delegatorAddr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return
	}
	validatorAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return
	}
	delegation, found := q.allianceKeeper.GetDelegation(ctx, delegatorAddr, validatorAddr, denom)
	if !found {
		return nil, alliancetypes.ErrDelegationNotFound
	}
	allianceValidator, err := q.allianceKeeper.GetAllianceValidator(ctx, validatorAddr)
	if err != nil {
		return nil, err
	}
	asset, found := q.allianceKeeper.GetAssetByDenom(ctx, denom)
	if !found {
		return nil, alliancetypes.ErrUnknownAsset
	}

	rewards, _, err := q.allianceKeeper.CalculateDelegationRewards(ctx, delegation, allianceValidator, asset)
	if err != nil {
		return
	}

	var coins []types.Coin
	for _, coin := range rewards {
		coins = append(coins, types.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		})
	}

	res, err = json.Marshal(types.DelegationRewardsResponse{
		Rewards: coins,
	})
	return res, err
}
