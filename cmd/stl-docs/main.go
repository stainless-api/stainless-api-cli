// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/stainless-api/stainless-api-cli/pkg/cmd"
	docs "github.com/urfave/cli-docs/v3"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "stl-docs",
		Usage: "Generate STL documentation in manpage and/or markdown formats",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "manpage",
				Usage: "generate manpage documentation `FILE`",
				Value: "stl.1",
			},
			&cli.StringFlag{
				Name:  "markdown",
				Usage: "generate markdown documentation `FILE`",
				Value: "usage.md",
			},
		},
		Action: generateDocs,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func generateDocs(ctx context.Context, c *cli.Command) error {
	manpageRequested := c.IsSet("manpage")
	markdownRequested := c.IsSet("markdown")

	if !manpageRequested && !markdownRequested {
		return cli.ShowAppHelp(c)
	}

	if manpageRequested {
		if err := generateManpage(c.String("manpage")); err != nil {
			return fmt.Errorf("failed to generate manpage: %w", err)
		}
	}

	if markdownRequested {
		if err := generateMarkdown(c.String("markdown")); err != nil {
			return fmt.Errorf("failed to generate markdown: %w", err)
		}
	}

	return nil
}

func generateManpage(filename string) error {
	fmt.Printf("Generating manpage: %s\n", filename)

	app := cmd.Command

	manpage, err := docs.ToMan(&app)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(manpage); err != nil {
		panic(err)
	}
	return nil
}

func generateMarkdown(filename string) error {
	fmt.Printf("Generating markdown: %s\n", filename)

	app := cmd.Command

	md, err := docs.ToMarkdown(&app)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString("# Stainless CLI\n\n" + md); err != nil {
		return err
	}
	return nil
}
