package tests_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/terra-money/alliance/x/alliance/types"
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
	indexKey := types.GetRedelegationIndexKey(srcValAddr, completion, denom, dstValAddr, delAddr)
	parsedDelKey, parsedTime, err := types.ParseRedelegationIndexForRedelegationKey(indexKey)
	require.NoError(t, err)
	require.Equal(t, parsedTime, completion)
	delKey := types.GetRedelegationKey(delAddr, denom, dstValAddr, completion)
	require.Equal(t, parsedDelKey, delKey)
}

func TestUndelegationIndex(t *testing.T) {
	delAddr, err := sdk.AccAddressFromHexUnsafe("aa")
	require.NoError(t, err)
	srcValAddr, err := sdk.ValAddressFromHex("bb")
	require.NoError(t, err)
	completion := time.Now().UTC()
	denom := "token"

	indexKey := types.GetUnbondingIndexKey(srcValAddr, completion, denom, delAddr)
	parsedUndelKey, parsedTime, err := types.ParseUnbondingIndexKeyToUndelegationKey(indexKey)
	require.NoError(t, err)
	require.Equal(t, parsedTime, completion)
	delKey := types.GetUndelegationQueueKey(completion, delAddr)
	require.Equal(t, delKey, parsedUndelKey)
}

func TestRewardWeightDecayQueueKey(t *testing.T) {
	triggerTime := time.Now().UTC()
	key := types.GetRewardWeightDecayQueueKey(triggerTime, "denom")
	parsedTime, denom := types.ParseRewardWeightDecayQueueKeyForDenom(key)
	require.Equal(t, "denom", denom)
	require.Equal(t, triggerTime, parsedTime)
}

func TestRewardSnapshotKey(t *testing.T) {
	denom := "denom"
	valAddr, err := sdk.ValAddressFromHex("bb")
	require.NoError(t, err)
	height := uint64(100)
	key := types.GetRewardWeightChangeSnapshotKey(denom, valAddr, height)

	parsedDenom, parsedValAddr, parsedHeight := types.ParseRewardWeightChangeSnapshotKey(key)
	require.Equal(t, denom, parsedDenom)
	require.Equal(t, valAddr, parsedValAddr)
	require.Equal(t, height, parsedHeight)
}

func TestValidatorKey(t *testing.T) {
	valAddr, err := sdk.ValAddressFromHex("bb")
	require.NoError(t, err)
	key := types.GetAllianceValidatorInfoKey(valAddr)

	parseValAddr := types.ParseAllianceValidatorKey(key)
	require.Equal(t, parseValAddr, valAddr)
}
