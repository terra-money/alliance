package tests_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/terra-money/alliance/x/alliance/types"
)

func TestSubtractDecCoinsWithRounding(t *testing.T) {
	// Normal case
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1000.00")),
	)

	c, err := types.SubtractDecCoinsWithRounding(a, b)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithSmallErrors(t *testing.T) {
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1000.90")),
	)

	c, err := types.SubtractDecCoinsWithRounding(a, b)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithBigErrors(t *testing.T) {
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1010.10")),
	)

	_, err := types.SubtractDecCoinsWithRounding(a, b)
	require.Error(t, err)
}
