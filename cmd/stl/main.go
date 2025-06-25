// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/stainless-api/stainless-api-cli/pkg/cmd"
	"github.com/stainless-api/stainless-api-go"
)

func main() {
	app := cmd.Command
	if err := app.Run(context.Background(), os.Args); err != nil {
		var apierr *stainless.Error
		if errors.As(err, &apierr) {
			fmt.Printf("%s\n", cmd.ColorizeJSON(apierr.RawJSON(), os.Stderr))
		} else {
			fmt.Printf("%s\n", err.Error())
		}
		os.Exit(1)
	}
}
