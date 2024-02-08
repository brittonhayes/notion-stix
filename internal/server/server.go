package server

import (
	"context"
	"fmt"
	"net/http"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/secure"
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
	ServiceName string
	Environment string
	Port        int
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

	// FIXME content security policy needs to be dialed in
	// csp := cspbuilder.Builder{
	// 	Directives: map[string][]string{
	// 		cspbuilder.DefaultSrc: {"self"},
	// 		cspbuilder.StyleSrc:   {"self", "https://cdn.tailwindcss.com", "https://unpkg.com/"},
	// 		cspbuilder.ScriptSrc:  {"self", "https://cdn.tailwindcss.com", "https://unpkg.com"},
	// 	},
	// }
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:       []string{"railway.app", "notion-stix.up.railway.app", "www.notion.so", "api.notion.com"},
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: csp.MustBuild(),
	})

	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(secureMiddleware.Handler)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Logger)

	api.Handler(config.Service, api.WithRouter(r))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(notionstix.HTML_HOME)
	})

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
