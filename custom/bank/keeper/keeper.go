package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	banktypes "github.com/terra-money/alliance/custom/bank/types"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Keeper struct {
	bankkeeper.BaseKeeper

	ak   alliancekeeper.Keeper
	sk   banktypes.StakingKeeper
	acck accountkeeper.AccountKeeper
}

var (
	_ bankkeeper.Keeper = Keeper{}
)

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak accountkeeper.AccountKeeper,
	paramSpace paramtypes.Subspace,
	blockedAddrs map[string]bool,
) Keeper {
	keeper := Keeper{
		BaseKeeper: bankkeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs),
		ak:         alliancekeeper.Keeper{},
		sk:         stakingkeeper.Keeper{},
		acck:       ak,
	}
	return keeper
}

func (k *Keeper) RegisterKeepers(ak alliancekeeper.Keeper, sk banktypes.StakingKeeper) {
	k.ak = ak
	k.sk = sk
}

// SupplyOf implements the Query/SupplyOf gRPC method
func (k Keeper) SupplyOf(c context.Context, req *types.QuerySupplyOfRequest) (*types.QuerySupplyOfResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	supply := k.GetSupply(ctx, req.Denom)

	if req.Denom == k.sk.BondDenom(ctx) {
		assets := k.ak.GetAllAssets(ctx)
		totalRewardWeights := sdk.ZeroDec()
		for _, asset := range assets {
			totalRewardWeights = totalRewardWeights.Add(asset.RewardWeight)
		}
		allianceBonded := k.ak.GetAllianceBondedAmount(ctx, k.acck.GetModuleAddress(alliancetypes.ModuleName))
		supply.Amount = supply.Amount.Sub(allianceBonded)
	}

	return &types.QuerySupplyOfResponse{Amount: sdk.NewCoin(req.Denom, supply.Amount)}, nil
}

// TotalSupply implements the Query/TotalSupply gRPC method
func (k Keeper) TotalSupply(ctx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalSupply, pageRes, err := k.GetPaginatedTotalSupply(sdkCtx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	allianceBonded := k.ak.GetAllianceBondedAmount(sdkCtx, k.acck.GetModuleAddress(alliancetypes.ModuleName))
	bondDenom := k.sk.BondDenom(sdkCtx)
	if totalSupply.AmountOf(bondDenom).IsPositive() {
		totalSupply = totalSupply.Sub(sdk.NewCoin(bondDenom, allianceBonded))
	}

	return &types.QueryTotalSupplyResponse{Supply: totalSupply, Pagination: pageRes}, nil
}
