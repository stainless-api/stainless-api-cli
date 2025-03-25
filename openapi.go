// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
  "context"
  "flag"
  "fmt"
  "net/http"
  "os"

  "github.com/stainless-sdks/stainless-v0-go"
  "github.com/stainless-sdks/stainless-v0-go/option"
)

func createOpenAPIRetrieveSubcommand() (Subcommand) {
  query := []byte("{}")
  header := []byte("{}")
  var flagSet = flag.NewFlagSet("openapi.retrieve", flag.ExitOnError)

  return Subcommand{
    flagSet: flagSet,
    handle: func(client *stainlessv0.Client) {
    res, err := client.OpenAPI.Get(context.TODO(), option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
      r.URL.RawQuery = serializeQuery(query).Encode()
      r.Header = serializeHeader(header)
      return mn(r)
    }))
    if err != nil {
      fmt.Printf("%s\n", err)
      os.Exit(1)
    }

    fmt.Printf("%s\n", res.JSON.RawJSON())
  },
  }
}
