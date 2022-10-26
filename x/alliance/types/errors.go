package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrEmptyValidatorAddr = sdkerrors.Register(ModuleName, 10, "empty validator address")
	ErrValidatorNotFound  = sdkerrors.Register(ModuleName, 11, "validator not found")

	ErrZeroDelegations = sdkerrors.Register(ModuleName, 20, "there are no delegations yet")

	ErrUnknownAsset = sdkerrors.Register(ModuleName, 30, "alliance asset is not whitelisted")
)
