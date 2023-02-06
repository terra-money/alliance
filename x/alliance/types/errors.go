package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrInvalidGenesisState = sdkerrors.Register(ModuleName, 0, "invalid genesis state")

	ErrEmptyValidatorAddr = sdkerrors.Register(ModuleName, 10, "empty validator address")
	ErrValidatorNotFound  = sdkerrors.Register(ModuleName, 11, "validator not found")

	ErrZeroDelegations    = sdkerrors.Register(ModuleName, 20, "there are no delegations yet")
	ErrInsufficientTokens = sdkerrors.Register(ModuleName, 21, "insufficient tokens")

	ErrUnknownAsset = sdkerrors.Register(ModuleName, 30, "alliance asset is not whitelisted")

	ErrRewardWeightOutOfBound = sdkerrors.Register(ModuleName, 40, "alliance asset must be between reward_weight_range")
)
