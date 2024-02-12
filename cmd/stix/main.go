// Package main is the entry point of the CLI application.
// It initializes the necessary components, sets up command-line flags,
// and executes the appropriate actions based on the user's input.
package main

import (
	"fmt"
	"os"
	"sort"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/kv"
	"github.com/brittonhayes/notion-stix/internal/mitre"
	"github.com/brittonhayes/notion-stix/internal/server"
	"github.com/brittonhayes/notion-stix/internal/service"
	"github.com/brittonhayes/notion-stix/internal/tasks"
	"github.com/charmbracelet/log"
	"github.com/hibiken/asynq"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	var (
		repo  notionstix.Repository
		store notionstix.Store
		queue *tasks.Queue
	)

	logger := log.New(os.Stdout)

	app := &cli.App{
		Name:                 "notion-stix",
		Authors:              []*cli.Author{{Name: "Britton Hayes"}},
		Usage:                "An integration for importing STIX-format Threat Intelligence into Notion",
		EnableBashCompletion: true,
		Args:                 false,
		Suggest:              true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "notion-auth-url",
				Aliases:  []string{"a"},
				Usage:    "The Notion auth URL (https://www.notion.so/my-integrations)",
				EnvVars:  []string{"NOTION_AUTH_URL"},
				Category: "Auth",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "redirect-uri",
				Usage:    "The redirect URI for the Notion OAuth integration",
				EnvVars:  []string{"REDIRECT_URI"},
				Required: true,
				Category: "Auth",
			},
			&cli.StringFlag{
				Name:     "client-id",
				Aliases:  []string{"i"},
				Usage:    "The Notion OAuth client ID",
				EnvVars:  []string{"OAUTH_CLIENT_ID"},
				Required: true,
				Category: "Auth",
			},
			&cli.StringFlag{
				Name:     "client-secret",
				Aliases:  []string{"s"},
				Usage:    "The Notion OAuth client secret",
				EnvVars:  []string{"OAUTH_CLIENT_SECRET"},
				Required: true,
				Category: "Auth",
			},
			&cli.StringFlag{
				Name:     "cookie-secret",
				Usage:    "The secret key used to encrypt cookies",
				EnvVars:  []string{"COOKIE_SECRET"},
				Required: true,
				Category: "Auth",
			},
			&cli.StringFlag{
				Name:     "db",
				Usage:    "The database to use for storing the STIX data",
				Value:    "notion-stix.db",
				EnvVars:  []string{"DB"},
				Category: "Application",
			},
			&cli.StringFlag{
				Name:     "redis-host",
				Usage:    "The host for the Redis server",
				Required: true,
				EnvVars:  []string{"REDISHOST"},
				Category: "Application",
			},
			&cli.IntFlag{
				Name:     "redis-port",
				Usage:    "The port for the Redis server",
				Required: true,
				EnvVars:  []string{"REDISPORT"},
				Category: "Application",
			},
			&cli.StringFlag{
				Name:     "redis-password",
				Usage:    "The password for the Redis server",
				EnvVars:  []string{"REDIS_PASSWORD"},
				Required: true,
				Category: "Application",
			},
			&cli.StringFlag{
				Name:     "page-id",
				Usage:    "The UUID of the Notion page to create the databases within",
				Category: "Application",
			},
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Usage:    "The port to run the server on",
				EnvVars:  []string{"PORT"},
				Value:    8080,
				Category: "Application",
			},
		},
		Before: func(c *cli.Context) error {
			b, err := notionstix.FS.ReadFile(mitre.STIX_JSON)
			if err != nil {
				return err
			}

			repo = mitre.NewRepository(b)

			// TODO at some point we need to do garbage collection
			// when? not sure yet. The docs mention doing it on nil errs
			// https://dgraph.io/docs/badger/get-started/#garbage-collection
			store, err = kv.NewPersistentKV(c.String("db"))
			if err != nil {
				return err
			}

			redisURL := fmt.Sprintf("%s:%d", c.String("redis-url"), c.Int("port"))
			queue = tasks.NewQueue(redisURL, c.String("redis-password"))

			return nil
		},
		Action: func(c *cli.Context) error {
			config := &server.Config{
				Repository: repo,
				Service: service.New(
					repo,
					c.String("redirect-uri"),
					c.String("client-id"),
					c.String("client-secret"),
					c.String("cookie-secret"),
					store,
					queue,
				),
				ServiceName: "stix",
				Environment: "production",
				Port:        c.Int("port"),
			}
			s := server.New(c.Context, config)

			redisURL := fmt.Sprintf("%s:%d", c.String("redis-url"), c.Int("port"))

			redisOpts := asynq.RedisClientOpt{Addr: redisURL, Password: c.String("redis-password"), DB: 0}

			g := new(errgroup.Group)
			g.Go(func() error {
				mux := tasks.NewMux()
				mux.Handle(tasks.TypeDatabaseCreate, tasks.NewAttackPatternProcessor())

				logger.Info("Starting queue server")
				queueServer := asynq.NewServer(redisOpts, asynq.Config{})
				return queueServer.Run(mux)
			})

			g.Go(func() error {
				logger.Info("Starting server", "port", config.Port, "service", config.ServiceName)
				return s.ListenAndServe()
			})

			// queue := asynq.NewClient(redisOpts)
			// defer queue.Close()

			return g.Wait()
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		logger.Error("Error running CLI", "error", err)
	}
}
