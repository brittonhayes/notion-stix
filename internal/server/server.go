package server

import (
	"context"
	"fmt"
	"net/http"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server represents an HTTP server.
type Server struct {
	server *http.Server
	Router *chi.Mux
}

// Config represents the configuration for the server.
type Config struct {
	Repository  notionstix.Repository
	Service     api.ServerInterface
	Addr        string
	ServiceName string
	Environment string
	Port        int
	Tracing     bool
}

// ListenAndServe starts the HTTP server and listens for incoming requests.
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
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

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Logger)

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
		server: httpsrv,
		Router: r,
	}
}
