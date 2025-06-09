// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-go/option"
	"github.com/urfave/cli/v3"
)

var orgsRetrieve = cli.Command{
	Name:  "retrieve",
	Usage: "Retrieve an organization by name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "org",
		},
	},
	Before:          initAPICommand,
	Action:          handleOrgsRetrieve,
	HideHelpCommand: true,
}

var orgsList = cli.Command{
	Name:            "list",
	Usage:           "List organizations the user has access to",
	Flags:           []cli.Flag{},
	Before:          initAPICommand,
	Action:          handleOrgsList,
	HideHelpCommand: true,
}

func handleOrgsRetrieve(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	res, err := cc.client.Orgs.Get(
		context.TODO(),
		cmd.Value("org").(string),
		option.WithMiddleware(cc.AsMiddleware()),
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}

func handleOrgsList(ctx context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(ctx, cmd)
	res, err := cc.client.Orgs.List(context.TODO(), option.WithMiddleware(cc.AsMiddleware()))
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", colorizeJSON(res.RawJSON(), os.Stdout))
	return nil
}
