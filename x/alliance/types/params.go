package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"time"
)

var (
	RewardDelayTime = []byte("RewardDelayTime")
)

var _ paramtypes.ParamSet = (*Params)(nil)

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(RewardDelayTime, &p.RewardDelayTime, validatePositiveDuration),
	}
}

func validatePositiveDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid paramter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: &d", v)
	}
	return nil
}
