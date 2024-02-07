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

	logger := log.New(os.Stdout)

	app := &cli.App{
		Name:  "notion-stix",
		Usage: "An integration for importing STIX-format Threat Intelligence into Notion",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "page-id",
				Usage: "The UUID of the Notion page to create the databases within",
			},
			&cli.StringFlag{
				Name:    "notion-auth-url",
				Aliases: []string{"a"},
				Usage:   "The Notion auth URL (https://www.notion.so/my-integrations)",
				EnvVars: []string{"NOTION_AUTH_URL"},
			},
			&cli.StringFlag{
				Name:    "redirect-uri",
				Usage:   "The redirect URI for the Notion OAuth integration",
				EnvVars: []string{"REDIRECT_URI"},
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
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "The port to run the server on",
				EnvVars: []string{"PORT"},
				Value:   8080,
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
			config := &server.Config{
				Repository:  repo,
				Service:     service.New(repo, c.String("redirect-uri"), c.String("client-id"), c.String("client-secret")),
				ServiceName: "stix",
				Environment: "production",
				Port:        c.Int("port"),
			}
			s := server.New(c.Context, config)

			logger.Info("Starting server", "port", config.Port, "service", config.ServiceName)
			return s.ListenAndServe()
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		logger.Error("Error running CLI", "error", err)
	}
}
