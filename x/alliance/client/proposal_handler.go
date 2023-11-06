package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/terra-money/alliance/x/alliance/client/cli"
)

var (
	CreateAllianceProposalHandler = govclient.NewProposalHandler(cli.CreateAlliance)
	UpdateAllianceProposalHandler = govclient.NewProposalHandler(cli.UpdateAlliance)
	DeleteAllianceProposalHandler = govclient.NewProposalHandler(cli.DeleteAlliance)
)
