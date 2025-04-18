// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
  "flag"
  "fmt"
  "log"
  "os"

  "github.com/stainless-sdks/stainless-v0-go"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Println("Expected subcommand")
    os.Exit(1)
  }

  subcommand := subcommands[os.Args[1]]
  if subcommand == nil {
    log.Fatalf("Unknown subcommand '%s'", os.Args[1])
  }

  subcommand.flagSet.Parse(os.Args[2:])

  var client *stainlessv0.Client = stainlessv0.NewClient()
  subcommand.handle(client)
}

func init() {
  initialBody := getStdInput()
  if initialBody == nil {
    initialBody = []byte("{}")
  }

  var openAPIRetrieveSubcommand = createOpenAPIRetrieveSubcommand()
  subcommands[openAPIRetrieveSubcommand.flagSet.Name()] = &openAPIRetrieveSubcommand

  var projectsUpdateSubcommand = createProjectsUpdateSubcommand(initialBody)
  subcommands[projectsUpdateSubcommand.flagSet.Name()] = &projectsUpdateSubcommand

  var projectsBranchesCreateSubcommand = createProjectsBranchesCreateSubcommand(initialBody)
  subcommands[projectsBranchesCreateSubcommand.flagSet.Name()] = &projectsBranchesCreateSubcommand

  var projectsBranchesRetrieveSubcommand = createProjectsBranchesRetrieveSubcommand()
  subcommands[projectsBranchesRetrieveSubcommand.flagSet.Name()] = &projectsBranchesRetrieveSubcommand

  var projectsSnippetsCreateRequestSubcommand = createProjectsSnippetsCreateRequestSubcommand(initialBody)
  subcommands[projectsSnippetsCreateRequestSubcommand.flagSet.Name()] = &projectsSnippetsCreateRequestSubcommand

  var buildsCreateSubcommand = createBuildsCreateSubcommand(initialBody)
  subcommands[buildsCreateSubcommand.flagSet.Name()] = &buildsCreateSubcommand

  var buildsRetrieveSubcommand = createBuildsRetrieveSubcommand()
  subcommands[buildsRetrieveSubcommand.flagSet.Name()] = &buildsRetrieveSubcommand

  var buildsListSubcommand = createBuildsListSubcommand()
  subcommands[buildsListSubcommand.flagSet.Name()] = &buildsListSubcommand

  var buildTargetOutputsListSubcommand = createBuildTargetOutputsListSubcommand()
  subcommands[buildTargetOutputsListSubcommand.flagSet.Name()] = &buildTargetOutputsListSubcommand

  var webhooksPostmanCreateNotificationSubcommand = createWebhooksPostmanCreateNotificationSubcommand(initialBody)
  subcommands[webhooksPostmanCreateNotificationSubcommand.flagSet.Name()] = &webhooksPostmanCreateNotificationSubcommand
}

var subcommands = map[string]*Subcommand{}

type Subcommand struct {
  flagSet *flag.FlagSet
  handle func(*stainlessv0.Client)
}
