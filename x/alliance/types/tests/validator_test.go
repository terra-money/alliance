package tests_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/alliance/x/alliance/types"
)

func TestSubtractDecCoinsWithRounding(t *testing.T) {
	// Normal case
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithSmallErrors(t *testing.T) {
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("1000.90")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithBigErrors(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("1010.10")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdkmath.LegacyMustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdkmath.LegacyMustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdkmath.LegacyMustNewDecFromStr("1000.00")),
	), c)
}
