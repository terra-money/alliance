package types_test

import (
	"alliance/x/alliance/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedelegationKey(t *testing.T) {
	delAddr, err := sdk.AccAddressFromHexUnsafe("aa")
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromHex("bb")
	require.NoError(t, err)
	completion := time.Now().UTC()
	key := types.GetRedelegationKey(delAddr, "denom", valAddr, completion)
	parsedCompletion := types.ParseRedelegationKeyForCompletionTime(key)
	require.Equal(t, completion, parsedCompletion)
}

func TestRedelegationQueueKey(t *testing.T) {
	completion := time.Now().UTC()
	key := types.GetRedelegationQueueKey(completion)
	parsedCompletion := types.ParseRedelegationQueueKey(key)
	require.Equal(t, completion, parsedCompletion)
}

func TestRedelegationIndex(t *testing.T) {
	delAddr, err := sdk.AccAddressFromHexUnsafe("aa")
	require.NoError(t, err)
	srcValAddr, err := sdk.ValAddressFromHex("bb")
	require.NoError(t, err)
	dstValAddr, err := sdk.ValAddressFromHex("bb")
	require.NoError(t, err)
	completion := time.Now().UTC()
	denom := "token"
	indexKey := types.GetRedelegationIndex(srcValAddr, completion, denom, dstValAddr, delAddr)
	parsedDelKey, parsedTime, err := types.ParseRedelegationIndexForRedelegationKey(indexKey)
	require.NoError(t, err)
	require.Equal(t, parsedTime, completion)
	delKey := types.GetRedelegationKey(delAddr, denom, dstValAddr, completion)
	require.Equal(t, parsedDelKey, delKey)
}
