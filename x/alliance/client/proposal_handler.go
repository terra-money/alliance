package client

import (
	"alliance/x/alliance/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	CreateAllianceHandler = govclient.NewProposalHandler(cli.CreateAlliance)
	UpdateAllianceHandler = govclient.NewProposalHandler(cli.UpdateAlliance)
	DeleteAllianceHandler = govclient.NewProposalHandler(cli.DeleteAlliance)
)
