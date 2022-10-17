package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "alliance"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_signletimemodule"
)

var (
	ModuleAccKey = []byte{0x01}

	AssetKey = []byte{0x11}

	DelegationKey = []byte{0x21}
)

func GetAssetKey(denom string) []byte {
	return append(AssetKey, address.MustLengthPrefix([]byte(denom))...)
}

func GetDelegationWithDenomKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string) []byte {
	return append(GetDelegationKey(delAddr, valAddr), address.MustLengthPrefix([]byte(denom))...)
}

// GetDelegationKey creates the key for delegator bond with validator for all denoms
func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

// GetDelegationsKey creates the prefix for a delegator for all validators
func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, address.MustLengthPrefix(delAddr)...)
}
