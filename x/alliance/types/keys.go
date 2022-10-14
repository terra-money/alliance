package types

import "github.com/cosmos/cosmos-sdk/types/address"

const (
	// ModuleName is the name of the staking module
	ModuleName = "alliance"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the staking module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the staking module
	RouterKey = ModuleName
)

var (
	AssetKey = []byte{0x11}
)

func GetAssetKey(denom string) []byte {
	return append(AssetKey, address.MustLengthPrefix([]byte(denom))...)
}
