package client

import (
	"github.com/terra-money/alliance/x/alliance/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	CreateAllianceProposalHandler = govclient.NewProposalHandler(cli.CreateAlliance)
	UpdateAllianceProposalHandler = govclient.NewProposalHandler(cli.UpdateAlliance)
	DeleteAllianceProposalHandler = govclient.NewProposalHandler(cli.DeleteAlliance)
)
