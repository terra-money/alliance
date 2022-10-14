package alliance_test

import (
	"testing"

	keepertest "alliance/testutil/keeper"
	"alliance/testutil/nullify"
	"alliance/x/alliance"
	"alliance/x/alliance/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AllianceKeeper(t)
	alliance.InitGenesis(ctx, *k, genesisState)
	got := alliance.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
