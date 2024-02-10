package service

import (
	"net/http"
	"os"

	"github.com/charmbracelet/log"

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
	repo notionstix.Repository

	logger *log.Logger

	client *http.Client

	redirectURI       string
	oauthClientID     string
	oauthClientSecret string
	cookieSecret      string

	// TODO consider replacing in-memory map with persistent token storage
	// Likely replace this with badger on-disk kv store and use a railway volume for persistence
	// https://dgraph.io/docs/badger/get-started/#using-keyvalue-pairs
	tokens map[string]string

	store notionstix.Store
}

// New creates a new instance of the Service.
func New(repo notionstix.Repository, redirectURI string, oauthClientID string, oauthClientSecret string, cookieSecret string, store notionstix.Store) *Service {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.Logger = nil
	return &Service{
		repo:              repo,
		tokens:            make(map[string]string),
		logger:            log.New(os.Stdout),
		client:            retryClient.StandardClient(),
		redirectURI:       redirectURI,
		oauthClientID:     oauthClientID,
		oauthClientSecret: oauthClientSecret,
		cookieSecret:      cookieSecret,
		store:             store,
	}
}
