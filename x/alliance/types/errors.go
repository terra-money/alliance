package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	ErrInvalidGenesisState = sdkerrors.Register(ModuleName, 0, "invalid genesis state")

	ErrEmptyValidatorAddr = sdkerrors.Register(ModuleName, 10, "empty validator address")
	ErrValidatorNotFound  = sdkerrors.Register(ModuleName, 11, "validator not found")
	ErrDelegationNotFound = sdkerrors.Register(ModuleName, 12, "delegation not found")

	ErrZeroDelegations    = sdkerrors.Register(ModuleName, 20, "there are no delegations yet")
	ErrInsufficientTokens = sdkerrors.Register(ModuleName, 21, "insufficient tokens")
	ErrAssetDissolving    = sdkerrors.Register(ModuleName, 22, "alliance operation not allowed because asset is dissolving")

	ErrUnknownAsset  = sdkerrors.Register(ModuleName, 30, "alliance asset is not whitelisted")
	ErrAlreadyExists = sdkerrors.Register(ModuleName, 31, "alliance asset already exists")

	ErrRewardWeightOutOfBound = sdkerrors.Register(ModuleName, 40, "alliance asset must be between reward_weight_range")
)
