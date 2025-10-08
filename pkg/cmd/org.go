// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/stainless-api/stainless-api-go/option"
	"github.com/tidwall/gjson"
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

func handleOrgsRetrieve(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("org") && len(unusedArgs) > 0 {
		cmd.Set("org", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
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

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("orgs retrieve", json, format, transform)
}

func handleOrgsList(_ context.Context, cmd *cli.Command) error {
	cc := getAPICommandContext(cmd)
	unusedArgs := cmd.Args().Slice()
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}
	var res []byte
	_, err := cc.client.Orgs.List(
		context.TODO(),
		option.WithMiddleware(cc.AsMiddleware()),
		option.WithResponseBodyInto(&res),
	)
	if err != nil {
		return err
	}

	json := gjson.Parse(string(res))
	format := cmd.Root().String("format")
	transform := cmd.Root().String("transform")
	return ShowJSON("orgs list", json, format, transform)
}
