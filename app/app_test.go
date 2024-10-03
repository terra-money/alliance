package app

import (
	"testing"

	"cosmossdk.io/log"
	db "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
)

func TestAppExportAndBlockedAddrs(t *testing.T) {
	app := Setup(t)
	_, err := app.ExportAppStateAndValidators(true, []string{}, nil)
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")

	app = New(
		log.NewNopLogger(),
		db.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		DefaultNodeHome,
		0,
		EmptyAppOptions{},
	)
	blockedAddrs := app.BlockedModuleAccountAddrs()

	require.NotContains(t, blockedAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func TestProtoAnnotations(t *testing.T) {
	r, err := proto.MergedRegistry()
	require.NoError(t, err)
	err = msgservice.ValidateProtoAnnotations(r)
	require.NoError(t, err)
}
