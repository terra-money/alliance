package alliance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gorilla/mux"
	"math/rand"

	// this line is used by starport scaffolding # 1

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/terra-money/alliance/x/alliance/client/cli"
	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
	_ module.AppModule      = AppModule{}
)

type AppModuleBasic struct {
	cdc  codec.Codec
	pcdc *codec.ProtoCodec
}

func (a AppModuleBasic) RegisterRESTRoutes(c client.Context, router *mux.Router) {
	// do nothing
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, sk types.StakingKeeper,
	ak types.AccountKeeper, bk types.BankKeeper, registry cdctypes.InterfaceRegistry) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc, pcdc: codec.NewProtoCodec(registry)},
		keeper:         keeper,
		stakingKeeper:  sk,
		bankKeeper:     bk,
		accountKeeper:  ak,
	}
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (a AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (a AppModuleBasic) DefaultGenesis(jsonCodec codec.JSONCodec) json.RawMessage {
	return jsonCodec.MustMarshalJSON(DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(jsonCodec codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	var genesis types.GenesisState
	if err := jsonCodec.UnmarshalJSON(message, &genesis); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return ValidateGenesis(&genesis)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule implements an application module for the alliance module.
type AppModule struct {
	AppModuleBasic
	keeper        keeper.Keeper
	stakingKeeper types.StakingKeeper
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
}

func (a AppModule) GenerateGenesisState(input *module.SimulationState) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) ProposalContents(simState module.SimulationState) []simulation.WeightedProposalContent {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RandomizedParams(r *rand.Rand) []simulation.ParamChange {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RegisterStoreDecoder(registry sdk.StoreDecoderRegistry) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) WeightedOperations(simState module.SimulationState) []simulation.WeightedOperation {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) BeginBlock(s sdk.Context, block abci.RequestBeginBlock) {
	// do nothing
}

func (a AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, a.keeper)
}

func (a AppModule) InitGenesis(ctx sdk.Context, jsonCodec codec.JSONCodec, message json.RawMessage) []abci.ValidatorUpdate {
	var genesis types.GenesisState
	jsonCodec.MustUnmarshalJSON(message, &genesis)
	return a.keeper.InitGenesis(ctx, &genesis)
}

func (a AppModule) ExportGenesis(ctx sdk.Context, jsonCodec codec.JSONCodec) json.RawMessage {
	genesis := a.keeper.ExportGenesis(ctx)
	return a.cdc.MustMarshalJSON(genesis)
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	RegisterInvariants(registry, a.keeper)
}

// Deprecated: use RegisterServices
func (a AppModule) Route() sdk.Route { return sdk.Route{} }

// Deprecated: use RegisterServices
func (AppModule) QuerierRoute() string { return types.RouterKey }

// Deprecated: use RegisterServices
func (a AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(a.keeper))
}

func (a AppModule) ConsensusVersion() uint64 {
	return 3
}
