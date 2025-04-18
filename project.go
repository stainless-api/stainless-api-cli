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

func createProjectsUpdateSubcommand(initialBody []byte) (Subcommand) {
  var projectName *string = nil
  query := []byte("{}")
  header := []byte("{}")
  body := initialBody
  var flagSet = flag.NewFlagSet("projects.update", flag.ExitOnError)

  flagSet.Func(
    "project-name",
    "",
    func(string string) error {
      projectName = &string
      return nil
    },
  )

  flagSet.Func(
    "display-name",
    "",
    func(string string) error {
      var jsonErr error
      body, jsonErr = jsonSet(body, "display_name", string)
      if jsonErr != nil {
        return jsonErr
      }
      return nil
    },
  )

  return Subcommand{
    flagSet: flagSet,
    handle: func(client *stainlessv0.Client) {
    res, err := client.Projects.Update(
      context.TODO(),
      *projectName,
      stainlessv0.ProjectUpdateParams{},
      option.WithMiddleware(func(r *http.Request, mn option.MiddlewareNext) (*http.Response, error) {
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
      }),
      option.WithRequestBody("application/json", body),
    )
    if err != nil {
      fmt.Printf("%s\n", err)
      os.Exit(1)
    }

    fmt.Printf("%s\n", res.JSON.RawJSON())
  },
  }
}
