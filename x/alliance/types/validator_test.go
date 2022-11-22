package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/terra-money/alliance/x/alliance/types"
	"testing"
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

	c := types.SubtractDecCoinsWithRounding(a, b)
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

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	), c)
}

func TestSubtractDecCoinsWithRoundingWithBigErrors(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	a := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1000.00")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	)
	b := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("400.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("1010.10")),
	)

	c := types.SubtractDecCoinsWithRounding(a, b)
	require.Equal(t, sdk.NewDecCoins(
		sdk.NewDecCoinFromDec("aaa", sdk.MustNewDecFromStr("600.00")),
		sdk.NewDecCoinFromDec("bbb", sdk.MustNewDecFromStr("0")),
		sdk.NewDecCoinFromDec("ccc", sdk.MustNewDecFromStr("1000.00")),
	), c)
}
