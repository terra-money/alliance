package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"time"
)

const (
	// ModuleName is the name of the staking module
	ModuleName = "alliance"

	// RewardsPoolName is the name of the module account for rewards
	RewardsPoolName = "alliance_rewards"

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

	AssetKey                    = []byte{0x11}
	ValidatorInfoKey            = []byte{0x12}
	AssetRebalanceQueueKey      = []byte{0x13} // TODO: find a better name. This represents an event to update weightages after changes to voting power
	RewardRateChangeSnapshotKey = []byte{0x14}

	DelegationKey        = []byte{0x21}
	RedelegationKey      = []byte{0x22}
	RedelegationQueueKey = []byte{0x23}
	UndelegationQueueKey = []byte{0x24}

	// Indexes for querying
	RedelegationByValidatorIndexKey = []byte{0x31}
)

func GetAssetKey(denom string) []byte {
	return append(AssetKey, address.MustLengthPrefix([]byte(denom))...)
}

// GetDelegationKey key is in the format of delegator|validator|denom
func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress, denom string) []byte {
	return append(GetDelegationsKeyForAllDenoms(delAddr, valAddr), address.MustLengthPrefix(CreateDenomAddressPrefix(denom))...)
}

// GetDelegationsKeyForAllDenoms creates the key for delegator bond with validator for all denoms
func GetDelegationsKeyForAllDenoms(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

// GetDelegationsKey creates the prefix for a delegator for all validators
func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, address.MustLengthPrefix(delAddr)...)
}

func GetRedelegationsKeyByDelegator(delAddr sdk.AccAddress) []byte {
	return append(RedelegationKey, address.MustLengthPrefix(delAddr)...)
}

func GetRedelegationsKeyByDelegatorAndDenom(delAddr sdk.AccAddress, denom string) []byte {
	return append(GetRedelegationsKeyByDelegator(delAddr), address.MustLengthPrefix(CreateDenomAddressPrefix(denom))...)
}

func GetRedelegationsKey(delAddr sdk.AccAddress, denom string, dstValAddr sdk.ValAddress) []byte {
	return append(GetRedelegationsKeyByDelegatorAndDenom(delAddr, denom), address.MustLengthPrefix(dstValAddr)...)
}

func GetRedelegationKey(delAddr sdk.AccAddress, denom string, dstValAddr sdk.ValAddress, completion time.Time) []byte {
	bz := sdk.FormatTimeBytes(completion)
	return append(GetRedelegationsKey(delAddr, denom, dstValAddr), bz...)
}

func GetRedelegationQueueKey(completion time.Time) []byte {
	bz := sdk.FormatTimeBytes(completion)
	return append(RedelegationQueueKey, bz...)
}

func GetRedelegationIndex(srcVal sdk.ValAddress, completion time.Time, denom string, dstVal sdk.ValAddress, delAddr sdk.AccAddress) []byte {
	key := append(GetRedelegationsIndexOrderedByValidatorKey(srcVal), address.MustLengthPrefix(sdk.FormatTimeBytes(completion))...)
	key = append(key, address.MustLengthPrefix(CreateDenomAddressPrefix(denom))...)
	key = append(key, address.MustLengthPrefix(dstVal)...)
	key = append(key, address.MustLengthPrefix(delAddr)...)
	return key
}

func GetRedelegationsIndexOrderedByValidatorKey(srcVal sdk.ValAddress) []byte {
	key := append(RedelegationByValidatorIndexKey, address.MustLengthPrefix(srcVal)...)
	return key
}

func ParseRedelegationIndexForRedelegationKey(key []byte) ([]byte, time.Time, error) {
	offset := 0
	offset += len(RedelegationByValidatorIndexKey)

	srcValAddrLen := int(key[offset])
	offset += 1
	offset += srcValAddrLen

	timeLen := int(key[offset])
	offset += 1
	timeBytes := key[offset : offset+timeLen]
	offset += timeLen

	denomLen := int(key[offset])
	offset += 1
	denomBytes := key[offset : offset+denomLen]
	offset += denomLen

	dstValAddrLen := int(key[offset])
	offset += 1
	dstValAddrBytes := key[offset : offset+dstValAddrLen]
	offset += dstValAddrLen

	delAddrLen := int(key[offset])
	offset += 1
	delAddrBytes := key[offset : offset+delAddrLen]

	newKey := append(RedelegationKey, address.MustLengthPrefix(delAddrBytes)...)
	newKey = append(newKey, address.MustLengthPrefix(denomBytes)...)
	newKey = append(newKey, address.MustLengthPrefix(dstValAddrBytes)...)
	newKey = append(newKey, timeBytes...)
	completionTime, err := sdk.ParseTimeBytes(timeBytes)
	return newKey, completionTime, err
}

func ParseRedelegationQueueKey(key []byte) time.Time {
	offset := 0
	offset += len(RedelegationQueueKey)
	b := key[offset:]
	t, err := sdk.ParseTimeBytes(b)
	if err != nil {
		panic(err)
	}
	return t
}

func CreateDenomAddressPrefix(denom string) []byte {
	// we add a "zero" byte at the end - null byte terminator, to allow prefix denom prefix
	// scan. Setting it is not needed (key[last] = 0) - because this is the default.
	key := make([]byte, len(denom)+1)
	copy(key, denom)
	return key
}

// ParseRedelegationKeyForCompletionTime key is in the format of RedelegationKey|delegator|denom|destination|timestamp
func ParseRedelegationKeyForCompletionTime(key []byte) time.Time {
	offset := 0
	offset += len(RedelegationKey)
	offset += int(key[offset]) + 1
	offset += int(key[offset]) + 1
	offset += int(key[offset]) + 1
	b := key[offset:]
	t, err := sdk.ParseTimeBytes(b)
	if err != nil {
		panic(err)
	}
	return t
}

func GetUndelegationQueueKey(completion time.Time) []byte {
	bz := sdk.FormatTimeBytes(completion)
	return append(UndelegationQueueKey, bz...)
}

func GetAllianceValidatorInfoKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorInfoKey, address.MustLengthPrefix(valAddr)...)
}

func ParseAllianceValidatorKey(key []byte) sdk.ValAddress {
	b := key[2:]
	return b
}

func GetRewardRateChangeSnapshotKey(denom string, val sdk.ValAddress, height uint64) []byte {
	key := append(RewardRateChangeSnapshotKey, address.MustLengthPrefix(CreateDenomAddressPrefix(denom))...)
	key = append(key, address.MustLengthPrefix(val)...)
	return append(key, sdk.Uint64ToBigEndian(height)...)
}
