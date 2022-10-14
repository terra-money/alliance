package alliance

import (
	"encoding/json"
	"fmt"
	"math/rand"

	// this line is used by starport scaffolding # 1

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"alliance/x/alliance/client/cli"
	"alliance/x/alliance/keeper"
	"alliance/x/alliance/simulation"
	"alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

var (
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModule           = AppModule{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModuleBasic struct {
	cdc codec.Codec
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, registry cdctypes.InterfaceRegistry) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
		stakingKeeper:  stakingKeeper,
	}
}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	//TODO implement me
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

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {
	//TODO implement me
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
}

func (a AppModule) InitGenesis(ctx sdk.Context, jsonCodec codec.JSONCodec, message json.RawMessage) []abci.ValidatorUpdate {
	var genesis types.GenesisState
	jsonCodec.MustUnmarshalJSON(message, &genesis)
	return a.keeper.InitGenesis(ctx, &genesis)
}

func (a AppModule) ExportGenesis(ctx sdk.Context, jsonCodec codec.JSONCodec) json.RawMessage {
	//TODO implement me
	return json.RawMessage{}
}

func (a AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {
	//TODO implement me
}

func (a AppModule) Route() sdk.Route {
	// Deprecated
	return sdk.Route{}
}

func (a AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (a AppModule) LegacyQuerierHandler(codec *codec.LegacyAmino) sdk.Querier {
	return keeper.NewLegacyQuerier(a.keeper, codec)
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (a AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(a.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), a.keeper)
}

func (a AppModule) ConsensusVersion() uint64 {
	return 3
}

func (a AppModule) GenerateGenesisState(input *module.SimulationState) {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (a AppModule) RandomizedParams(r *rand.Rand) []simtypes.ParamChange {
	//TODO implement me
	panic("implement me")
}

func (a AppModule) RegisterStoreDecoder(registry sdk.StoreDecoderRegistry) {
	registry[types.StoreKey] = simulation.NewDecodeStore(a.cdc)
}

func (a AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	//TODO implement me
	panic("implement me")
}
