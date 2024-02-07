// Package main is the entry point of the CLI application.
// It initializes the necessary components, sets up command-line flags,
// and executes the appropriate actions based on the user's input.
package main

import (
	"os"
	"sort"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/mitre"
	"github.com/brittonhayes/notion-stix/internal/server"
	"github.com/brittonhayes/notion-stix/internal/service"
	"github.com/charmbracelet/log"

	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// TODO interact with STIX data from cli with no Notion integration

func main() {

	var (
		repo notionstix.Repository
	)

	app := &cli.App{
		Name:  "notion-stix",
		Usage: "An integration for importing STIX-format Threat Intelligence into Notion",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "page-id",
				Aliases: []string{"p"},
				Usage:   "The UUID of the Notion page to create the databases within",
			},
			&cli.StringFlag{
				Name:    "notion-auth-url",
				Aliases: []string{"a"},
				Usage:   "The Notion auth URL (https://www.notion.so/my-integrations)",
				EnvVars: []string{"NOTION_AUTH_URL"},
			},
			&cli.StringFlag{
				Name:    "client-id",
				Aliases: []string{"i"},
				Usage:   "The Notion OAuth client ID",
				EnvVars: []string{"OAUTH_CLIENT_ID"},
			},
			&cli.StringFlag{
				Name:    "client-secret",
				Aliases: []string{"s"},
				Usage:   "The Notion OAuth client secret",
				EnvVars: []string{"OAUTH_CLIENT_SECRET"},
			},
		},
		Before: func(c *cli.Context) error {
			b, err := notionstix.FS.ReadFile(notionstix.MitreEnterpriseAttack.String())
			if err != nil {
				return err
			}

			repo = mitre.NewRepository(b)
			return nil
		},
		Action: func(c *cli.Context) error {
			s := server.New(c.Context, &server.Config{
				Repository:  repo,
				Service:     service.New(repo),
				ServiceName: "stix",
				Environment: "production",
				Port:        8080,
			})
			return s.ListenAndServe()
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Error("Error running CLI", "error", err)
	}
}
