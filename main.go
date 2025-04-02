// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
  "flag"
  "fmt"
  "log"
  "os"

  "github.com/stainless-api/stainless-api-go"
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

  var projectsConfigCreateBranchSubcommand = createProjectsConfigCreateBranchSubcommand(initialBody)
  subcommands[projectsConfigCreateBranchSubcommand.flagSet.Name()] = &projectsConfigCreateBranchSubcommand

  var projectsConfigCreateCommitSubcommand = createProjectsConfigCreateCommitSubcommand(initialBody)
  subcommands[projectsConfigCreateCommitSubcommand.flagSet.Name()] = &projectsConfigCreateCommitSubcommand

  var projectsConfigMergeSubcommand = createProjectsConfigMergeSubcommand(initialBody)
  subcommands[projectsConfigMergeSubcommand.flagSet.Name()] = &projectsConfigMergeSubcommand

  var buildsRetrieveSubcommand = createBuildsRetrieveSubcommand()
  subcommands[buildsRetrieveSubcommand.flagSet.Name()] = &buildsRetrieveSubcommand

  var buildsTargetRetrieveSubcommand = createBuildsTargetRetrieveSubcommand()
  subcommands[buildsTargetRetrieveSubcommand.flagSet.Name()] = &buildsTargetRetrieveSubcommand

  var buildsTargetArtifactsRetrieveSourceSubcommand = createBuildsTargetArtifactsRetrieveSourceSubcommand()
  subcommands[buildsTargetArtifactsRetrieveSourceSubcommand.flagSet.Name()] = &buildsTargetArtifactsRetrieveSourceSubcommand
}

var subcommands = map[string]*Subcommand{}

type Subcommand struct {
  flagSet *flag.FlagSet
  handle func(*stainlessv0.Client)
}
