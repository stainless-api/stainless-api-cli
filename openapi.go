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
      q := r.URL.Query()
      for key, values := range serializeQuery(query) {
        for _, value := range values {
          q.Add(key, value)
        }
      }
      r.URL.RawQuery = q.Encode()

      for key, values := range serializeHeader(header) {
        for _, value := range values {
          r.Header.Add(key, value)
        }
      }

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
