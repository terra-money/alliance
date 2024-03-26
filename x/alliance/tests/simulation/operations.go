package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/auth/tx"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(cdc *codec.ProtoCodec,
	ak types.AccountKeeper, bk types.BankKeeper,
	sk types.StakingKeeper, k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgDelegate     int
		weightMsgUndelegate   int
		weightMsgRedelegate   int
		weightMsgClaimRewards int
	)

	weightMsgDelegate = 100
	weightMsgUndelegate = 1
	weightMsgRedelegate = 20
	weightMsgClaimRewards = 10

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgDelegate,
			SimulateMsgDelegate(cdc, ak, bk, sk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRedelegate,
			SimulateMsgRedelegate(cdc, ak, bk, sk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgUndelegate,
			SimulateMsgUndelegate(cdc, ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgClaimRewards,
			SimulateMsgClaimRewards(cdc, ak, bk, k),
		),
	}
}

func SimulateMsgDelegate(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainId string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		assets := k.GetAllAssets(ctx)
		if len(assets) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "No assets"), nil, nil
		}
		idx := simtypes.RandIntBetween(r, 0, len(assets)-1)
		assetToDelegate := assets[idx]
		amountToDelegate := simtypes.RandomAmount(r, sdkmath.NewInt(1000_000_000))
		if amountToDelegate.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "0 delegate amount"), nil, nil
		}
		validators := sk.GetAllValidators(ctx)
		idx = simtypes.RandIntBetween(r, 0, len(validators)-1)
		validatorToDelegateTo := validators[idx]
		coinToDelegate := sdk.NewCoin(assetToDelegate.Denom, amountToDelegate)

		bk.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(coinToDelegate))                                        //nolint:errcheck
		bk.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, simAccount.Address, sdk.NewCoins(coinToDelegate)) //nolint:errcheck

		msg := &types.MsgDelegate{
			DelegatorAddress: simAccount.Address.String(),
			ValidatorAddress: validatorToDelegateTo.GetOperator(),
			Amount:           coinToDelegate,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           tx.NewTxConfig(cdc, tx.DefaultSignModes),
			Cdc:             cdc,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(coinToDelegate),
			Bankkeeper:      bk,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgRedelegate(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, sk types.StakingKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainId string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var delegations []types.Delegation
		k.IterateDelegations(ctx, func(d types.Delegation) bool {
			delegations = append(delegations, d)
			return false
		})
		if len(delegations) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "No delegations yet"), nil, nil
		}
		var idx int
		if len(delegations) == 1 {
			idx = 0
		} else {
			idx = simtypes.RandIntBetween(r, 0, len(delegations)-1)
		}
		delegation := delegations[idx]

		simAccountAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		simAccount, found := simtypes.FindAccount(accs, simAccountAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "Account not found"), nil, nil
		}
		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "Validator not found"), nil, nil
		}

		if k.HasRedelegation(ctx, simAccountAddr, valAddr, delegation.Denom) {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "Cannot perform redelegation from a previous destination of a redelegation"), nil, nil
		}

		asset, _ := k.GetAssetByDenom(ctx, delegation.Denom)
		bondedTokens := types.GetDelegationTokens(delegation, validator, asset)

		amountToRedelegate := simtypes.RandomAmount(r, bondedTokens.Amount)
		if amountToRedelegate.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "0 redelegate amount"), nil, nil
		}

		validators := sk.GetAllValidators(ctx)
		idx = simtypes.RandIntBetween(r, 0, len(validators)-1)
		validatorToDelegateTo := validators[idx]

		if delegation.ValidatorAddress == validatorToDelegateTo.GetOperator() {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "redelegation to the same validator"), nil, nil
		}

		msg := &types.MsgRedelegate{
			DelegatorAddress:    delegation.DelegatorAddress,
			ValidatorSrcAddress: delegation.ValidatorAddress,
			ValidatorDstAddress: validatorToDelegateTo.GetOperator(),
			Amount:              sdk.NewCoin(asset.Denom, amountToRedelegate),
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         tx.NewTxConfig(cdc, tx.DefaultSignModes),
			Cdc:           cdc,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    types.ModuleName,
			Bankkeeper:    bk,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgUndelegate(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainId string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var delegations []types.Delegation
		k.IterateDelegations(ctx, func(d types.Delegation) bool {
			delegations = append(delegations, d)
			return false
		})
		if len(delegations) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "No delegations yet"), nil, nil
		}
		var idx int
		if len(delegations) == 1 {
			idx = 0
		} else {
			idx = simtypes.RandIntBetween(r, 0, len(delegations)-1)
		}
		delegation := delegations[idx]

		simAccountAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		simAccount, found := simtypes.FindAccount(accs, simAccountAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "Account not found"), nil, nil
		}
		valAddr, _ := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		validator, err := k.GetAllianceValidator(ctx, valAddr)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "Validator not found"), nil, nil
		}

		asset, _ := k.GetAssetByDenom(ctx, delegation.Denom)
		bondedTokens := types.GetDelegationTokens(delegation, validator, asset)

		amountToUndelegate := simtypes.RandomAmount(r, bondedTokens.Amount.Sub(sdkmath.NewInt(1))).Add(sdkmath.NewInt(1))
		if amountToUndelegate.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "0 undelegate amount"), nil, nil
		}

		msg := &types.MsgUndelegate{
			DelegatorAddress: simAccount.Address.String(),
			ValidatorAddress: delegation.ValidatorAddress,
			Amount:           sdk.NewCoin(delegation.Denom, amountToUndelegate),
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         tx.NewTxConfig(cdc, tx.DefaultSignModes),
			Cdc:           cdc,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    types.ModuleName,
			Bankkeeper:    bk,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

func SimulateMsgClaimRewards(cdc *codec.ProtoCodec, ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainId string) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var delegations []types.Delegation
		k.IterateDelegations(ctx, func(d types.Delegation) bool {
			delegations = append(delegations, d)
			return false
		})
		if len(delegations) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "No delegations yet"), nil, nil
		}
		var idx int
		if len(delegations) == 1 {
			idx = 0
		} else {
			idx = simtypes.RandIntBetween(r, 0, len(delegations)-1)
		}
		delegation := delegations[idx]

		simAccountAddr, _ := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		simAccount, found := simtypes.FindAccount(accs, simAccountAddr)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.MsgRedelegateType, "Account not found"), nil, nil
		}

		msg := &types.MsgClaimDelegationRewards{
			DelegatorAddress: simAccount.Address.String(),
			ValidatorAddress: delegation.ValidatorAddress,
			Denom:            delegation.Denom,
		}

		txCtx := simulation.OperationInput{
			R:             r,
			App:           app,
			TxGen:         tx.NewTxConfig(cdc, tx.DefaultSignModes),
			Cdc:           cdc,
			Msg:           msg,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    types.ModuleName,
			Bankkeeper:    bk,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
