package tests_test

import (
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/alliance/x/alliance/types"
)

func TestSubtractDecCoinsWithRounding(t *testing.T) {
	// Normal case
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", math.LegacyMustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("1000.00")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", math.LegacyMustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithSmallErrors(t *testing.T) {
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", math.LegacyMustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("1000.90")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", math.LegacyMustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithBigErrors(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", math.LegacyMustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("1010.10")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", math.LegacyMustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", math.LegacyMustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", math.LegacyMustNewDecFromStr("1000.00")),
	), c)
}
