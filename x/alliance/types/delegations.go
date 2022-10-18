package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ReduceShares
func (d *Delegation) ReduceShares(shares sdk.Dec) {
	if d.Shares.LTE(shares) {
		d.Shares = sdk.ZeroDec()
	} else {
		d.Shares = d.Shares.Sub(shares)
	}
}

// AddShares
func (d *Delegation) AddShares(shares sdk.Dec) {
	d.Shares = d.Shares.Add(shares)
}
