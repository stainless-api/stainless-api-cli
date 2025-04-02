// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
  "context"
  "flag"
  "fmt"
  "net/http"
  "os"

  "github.com/stainless-api/stainless-api-go"
  "github.com/stainless-api/stainless-api-go/option"
)

func createBuildsRetrieveSubcommand() (Subcommand) {
  var buildID *string = nil
  query := []byte("{}")
  header := []byte("{}")
  var flagSet = flag.NewFlagSet("builds.retrieve", flag.ExitOnError)

  flagSet.Func(
    "build-id",
    "",
    func(string string) error {
      buildID = &string
      return nil
    },
  )

  return Subcommand{
    flagSet: flagSet,
    handle: func(client *stainlessv0.Client) {
    res, err := client.Builds.Get(
      context.TODO(),
      *buildID,
      option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
        r.URL.RawQuery = serializeQuery(query).Encode()
        r.Header = serializeHeader(header)
        return mn(r)
      }),
    )
    if err != nil {
      fmt.Printf("%s\n", err)
      os.Exit(1)
    }

    fmt.Printf("%s\n", res.JSON.RawJSON())
  },
  }
}
