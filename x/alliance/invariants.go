package alliance

import (
	"alliance/x/alliance/keeper"
	"alliance/x/alliance/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k keeper.Keeper) {
	ir.RegisterRoute(types.ModuleName, "validator-shares", ValidatorSharesInvariant(k))
	ir.RegisterRoute(types.ModuleName, "delegator-shares", DelegatorSharesInvariant(k))
}

func RunAllInvariants(ctx sdk.Context, k keeper.Keeper) (res string, stop bool) {
	res, stop = ValidatorSharesInvariant(k)(ctx)
	if stop {
		return res, stop
	}
	res, stop = DelegatorSharesInvariant(k)(ctx)
	return res, stop
}

func ValidatorSharesInvariant(k keeper.Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg    string
			broken bool
		)
		assets := k.GetAllAssets(ctx)
		infos := k.GetAllAllianceValidatorInfo(ctx)
		validatorShares := map[string]sdk.Dec{} // {denom: shares}
		for _, info := range infos {
			for _, share := range info.ValidatorShares {
				if validatorShares[share.Denom].IsNil() {
					validatorShares[share.Denom] = share.Amount
				} else {
					validatorShares[share.Denom] = validatorShares[share.Denom].Add(share.Amount)
				}
			}
		}
		for _, asset := range assets {
			if !validatorShares[asset.Denom].IsNil() && !asset.TotalValidatorShares.Equal(validatorShares[asset.Denom]) {
				broken = true
				msg += fmt.Sprintf("broken alliance validator share invariance: \n"+
					"asset.TotalValidatorShares: %s\n"+
					"sum of validator shares: %s\n", asset.TotalValidatorShares, validatorShares[asset.Denom])
			}
		}
		return sdk.FormatInvariant(types.ModuleName, "validator shares", msg), broken
	}
}

func DelegatorSharesInvariant(k keeper.Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg    string
			broken bool
		)
		delegatorShares := map[string]map[string]sdk.Dec{} // {validator: {asset: share}}
		k.IterateDelegations(ctx, func(delegation types.Delegation) bool {
			if delegatorShares[delegation.ValidatorAddress] == nil {
				delegatorShares[delegation.ValidatorAddress] = map[string]sdk.Dec{
					delegation.Denom: delegation.Shares,
				}
			} else {
				if delegatorShares[delegation.ValidatorAddress][delegation.Denom].IsNil() {
					delegatorShares[delegation.ValidatorAddress][delegation.Denom] = delegation.Shares
				} else {
					delegatorShares[delegation.ValidatorAddress][delegation.Denom] = delegatorShares[delegation.ValidatorAddress][delegation.Denom].Add(delegation.Shares)
				}
			}
			return false
		})

		for val, assets := range delegatorShares {
			valAddr, err := sdk.ValAddressFromBech32(val)
			if err != nil {
				msg = fmt.Sprintf("alliance validator address invalid\n")
				broken = true
				break
			}
			info, found := k.GetAllianceValidatorInfo(ctx, valAddr)
			if !found {
				msg = fmt.Sprintf("alliance validator info for %s not found\n", val)
				broken = true
				break
			}
			shares := sdk.NewDecCoins(info.TotalDelegatorShares...)
			for denom, amount := range assets {
				if !shares.AmountOf(denom).Equal(amount) {
					msg += fmt.Sprintf("broken alliance delegation share invariance: \n"+
						"validator.TotalDelegatorShares(%s): %s\n"+
						"sum of delegator shares: %s\n", denom, shares, amount)
					broken = true
				}
			}
		}
		return sdk.FormatInvariant(types.ModuleName, "validator shares", msg), broken
	}
}
