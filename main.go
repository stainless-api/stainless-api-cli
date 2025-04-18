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

	var projectsRetrieveSubcommand = createProjectsRetrieveSubcommand()
	subcommands[projectsRetrieveSubcommand.flagSet.Name()] = &projectsRetrieveSubcommand

	var projectsUpdateSubcommand = createProjectsUpdateSubcommand(initialBody)
	subcommands[projectsUpdateSubcommand.flagSet.Name()] = &projectsUpdateSubcommand

	var projectsBranchesCreateSubcommand = createProjectsBranchesCreateSubcommand(initialBody)
	subcommands[projectsBranchesCreateSubcommand.flagSet.Name()] = &projectsBranchesCreateSubcommand

	var projectsBranchesRetrieveSubcommand = createProjectsBranchesRetrieveSubcommand()
	subcommands[projectsBranchesRetrieveSubcommand.flagSet.Name()] = &projectsBranchesRetrieveSubcommand

	var buildsCreateSubcommand = createBuildsCreateSubcommand(initialBody)
	subcommands[buildsCreateSubcommand.flagSet.Name()] = &buildsCreateSubcommand

	var buildsRetrieveSubcommand = createBuildsRetrieveSubcommand()
	subcommands[buildsRetrieveSubcommand.flagSet.Name()] = &buildsRetrieveSubcommand

	var buildsListSubcommand = createBuildsListSubcommand()
	subcommands[buildsListSubcommand.flagSet.Name()] = &buildsListSubcommand

	var buildTargetOutputsRetrieveSubcommand = createBuildTargetOutputsRetrieveSubcommand()
	subcommands[buildTargetOutputsRetrieveSubcommand.flagSet.Name()] = &buildTargetOutputsRetrieveSubcommand
}

var subcommands = map[string]*Subcommand{}

type Subcommand struct {
	flagSet *flag.FlagSet
	handle  func(*stainlessv0.Client)
}
