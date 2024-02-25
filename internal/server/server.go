package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/httprate"
	"github.com/unrolled/secure"
)

// Server represents an HTTP server.
type Server struct {
	httpServer *http.Server
	Router     *chi.Mux
}

// Config represents the configuration for the server.
type Config struct {
	Repository  notionstix.Repository
	Service     api.ServerInterface
	ServiceName string
	Environment string
	RedisAddr   string
	Port        int
}

// ListenAndServe starts the HTTP server and listens for incoming requests.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// New creates a new instance of the server.
func New(ctx context.Context, config *Config) *Server {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	logger := httplog.NewLogger(config.ServiceName, httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelInfo,
		MessageFieldName: "message",
		QuietDownRoutes: []string{
			"/",
			"/healthz",
			"/api/events",
		},
		QuietDownPeriod: 10 * time.Second,
	})

	r := chi.NewRouter()

	// TODO content security policy needs to be dialed in
	// csp := cspbuilder.Builder{
	// 	Directives: map[string][]string{
	// 		cspbuilder.DefaultSrc: {"self"},
	// 		cspbuilder.StyleSrc:   {"self", "https://cdn.tailwindcss.com", "https://unpkg.com/"},
	// 		cspbuilder.ScriptSrc:  {"self", "https://cdn.tailwindcss.com", "https://unpkg.com"},
	// 	},
	// }
	allowedHosts := []string{"railway.app", "notion-stix.up.railway.app", "www.notion.so", "api.notion.com"}
	if config.Environment == "development" {
		allowedHosts = []string{}
	}

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:       allowedHosts,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: csp.MustBuild(),
	})

	r.Use(middleware.Heartbeat("/healthz"))

	// TODO: Dial in this rate limiter so it isn't so strict
	r.Use(httprate.Limit(10, 1*time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint, func(r *http.Request) (string, error) {
			cookie, err := r.Cookie("bot_id")
			if err != nil {
				return "", nil
			}
			return cookie.Value, nil
		}),
	))
	r.Use(secureMiddleware.Handler)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(httplog.RequestLogger(logger))

	api.Handler(config.Service, api.WithRouter(r))

	port := config.Port
	if port == 0 {
		port = 8080
	}

	httpsrv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
	}

	return &Server{
		httpServer: httpsrv,
		Router:     r,
	}
}
