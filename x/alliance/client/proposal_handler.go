package client

import (
	"github.com/terra-money/alliance/x/alliance/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

var (
	CreateAllianceProposalHandler = govclient.NewProposalHandler(cli.CreateAlliance, nil)
	UpdateAllianceProposalHandler = govclient.NewProposalHandler(cli.UpdateAlliance, nil)
	DeleteAllianceProposalHandler = govclient.NewProposalHandler(cli.DeleteAlliance, nil)
)
