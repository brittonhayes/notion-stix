package service

import (
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/time/rate"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	NOTION_URL       = "https://www.notion.so/"
	NOTION_OAUTH_URL = "https://api.notion.com/v1/oauth/token"

	ErrCancel       = "internal server error caused by user cancellation"
	ErrOAuthGrant   = "internal server error caused by oauth grant content"
	ErrMissingToken = "internal server error caused by missing oauth token"
	ErrTokenRequest = "internal server error caused by oauth request to Notion API"
	ErrTokenDecode  = "internal server error caused by decoding oauth token response"
	ErrImportSTIX   = "internal server error caused by importing STIX data to Notion"
)

// Service represents a service that handles integration setup and other operations.
type Service struct {
	repo    notionstix.Repository
	store   notionstix.Store
	limiter *rate.Limiter
	updates map[string]chan string

	client    *http.Client
	logger    *log.Logger
	templates *template.Template

	redirectURI       string
	oauthClientID     string
	oauthClientSecret string
	cookieSecret      string
}

// New creates a new instance of the Service.
func New(repo notionstix.Repository, redirectURI string, oauthClientID string, oauthClientSecret string, cookieSecret string, store notionstix.Store) *Service {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.Logger = nil
	templates := template.Must(template.ParseFS(notionstix.TEMPLATES, "web/*.html"))
	return &Service{
		repo:              repo,
		updates:           make(map[string]chan string),
		logger:            log.New(os.Stdout),
		templates:         templates,
		client:            retryClient.StandardClient(),
		redirectURI:       redirectURI,
		oauthClientID:     oauthClientID,
		oauthClientSecret: oauthClientSecret,
		cookieSecret:      cookieSecret,
		store:             store,
		limiter:           rate.NewLimiter(rate.Every(time.Second), 3),
	}
}
