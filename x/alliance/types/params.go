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
	RewardClaimInterval = []byte("RewardClaimInterval")
	LastRewardClaimTime = []byte("LastRewardClaimTime")
	GlobalRewardIndices = []byte("GlobalRewardIndices")
)

var _ paramtypes.ParamSet = (*Params)(nil)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	idxs := NewRewardIndices(p.GlobalRewardIndices)
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(RewardDelayTime, &p.RewardDelayTime, validatePositiveDuration),
		paramtypes.NewParamSetPair(RewardClaimInterval, &p.RewardClaimInterval, validatePositiveDuration),
		paramtypes.NewParamSetPair(LastRewardClaimTime, &p.LastRewardClaimTime, validateTime),
		paramtypes.NewParamSetPair(GlobalRewardIndices, &idxs, validatePositiveRewardIndices),
	}
}

func validatePositiveDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("duration must be positive: %d", v)
	}
	return nil
}

func validateTime(i interface{}) error {
	_, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
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
		RewardClaimInterval: time.Minute * 5,
		LastRewardClaimTime: time.Now(),
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
