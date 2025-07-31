// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/stainless-api/stainless-api-cli/pkg/cmd"
	"github.com/stainless-api/stainless-api-go"
)

func main() {
	app := cmd.Command
	if err := app.Run(context.Background(), os.Args); err != nil {
		var apierr *stainless.Error
		if errors.As(err, &apierr) {
			fmt.Fprintf(os.Stderr, "%s %q: %d %s\n", apierr.Request.Method, apierr.Request.URL, apierr.Response.StatusCode, http.StatusText(apierr.Response.StatusCode))
			fmt.Fprintf(os.Stdout, "%s\n", cmd.ColorizeJSON(apierr.RawJSON(), os.Stdout))
		} else {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		os.Exit(1)
	}
}
