package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	RewardDelayTime     = []byte("RewardDelayTime")
	GlobalRewardIndices = []byte("GlobalRewardIndices")
)

var _ paramtypes.ParamSet = (*Params)(nil)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	idxs := NewRewardIndices(p.GlobalRewardIndices)
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(RewardDelayTime, &p.RewardDelayTime, validatePositiveDuration),
		paramtypes.NewParamSetPair(GlobalRewardIndices, &idxs, validatePositiveRewardIndices),
	}
}

func validatePositiveDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}
	return nil
}

func validatePositiveRewardIndices(i interface{}) error {
	v, ok := i.(RewardIndices)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	for _, i := range v {
		if i.Index.LT(sdk.ZeroDec()) {
			return fmt.Errorf("unbonding time must be positive: %d", v)
		}
	}
	return nil
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		RewardDelayTime:     time.Hour,
		GlobalRewardIndices: make([]RewardIndex, 0),
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

type RewardIndices []RewardIndex

func NewRewardIndices(r []RewardIndex) RewardIndices {
	return r
}

func (r RewardIndices) GetIndexByDenom(denom string) (ri *RewardIndex, found bool) {
	idx := slices.IndexFunc(r, func(e RewardIndex) bool {
		return e.Denom == denom
	})
	if idx < 0 {
		return &RewardIndex{}, false
	} else {
		return &r[idx], true
	}
}
