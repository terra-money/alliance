package simulation

import (
	"bytes"
	"fmt"

	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.RewardDelayTime):
			var powerA, powerB sdk.IntProto

			cdc.MustUnmarshal(kvA.Value, &powerA)
			cdc.MustUnmarshal(kvB.Value, &powerB)

			return fmt.Sprintf("%v\n%v", powerA, powerB)
		default:
			panic(fmt.Sprintf("invalid key prefix %X", kvA.Key[:1]))
		}
	}
}
