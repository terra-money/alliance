package types

import (
	"context"
)

type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
}
