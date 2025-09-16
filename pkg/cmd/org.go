// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"

	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var orgsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve an organization by name.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "org",
		},
	},
	Action:          handleOrgsRetrieve,
	HideHelpCommand: true,
}

var orgsList = cli.Command{
	Name:            "list",
	Usage:           "List organizations accessible to the current authentication method.",
	Flags:           []cli.Flag{},
	Action:          handleOrgsList,
	HideHelpCommand: true,
}

func handleOrgsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	var res []byte
	_, err := cc.client.Orgs.Get(
		context.TODO(),
		cmd.Value("org").(string),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("orgs retrieve", string(res), format)
}

func handleOrgsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	var res []byte
	_, err := cc.client.Orgs.List(
		context.TODO(),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	format := cmd.Root().String("format")
	return ShowJSON("orgs list", string(res), format)
}
