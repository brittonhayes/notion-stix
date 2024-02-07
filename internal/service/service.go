package service

import (
	"fmt"
	"net/http"
	"net/url"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	NOTION_OAUTH_URL = "https://api.notion.com/v1/oauth/token"

	ErrOAuthGrant   = "internal server error caused by oauth grant content"
	ErrTokenRequest = "internal server error caused by oauth request to Notion API"
	ErrTokenDecode  = "internal server error caused by decoding oauth token response"
	ErrImportSTIX   = "internal server error caused by importing STIX data to Notion"
)

// Service represents a service that handles integration setup and other operations.
type Service struct {
	repo   notionstix.Repository
	source notionstix.StixSource

	client *http.Client

	redirectURI       string
	oauthClientID     string
	oauthClientSecret string
}

// New creates a new instance of the Service.
func New(repo notionstix.Repository, redirectURI string, oauthClientID string, oauthClientSecret string) Service {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.Logger = nil
	return Service{repo: repo, client: retryClient.StandardClient(), redirectURI: redirectURI, oauthClientID: oauthClientID, oauthClientSecret: oauthClientSecret}
}

func (s Service) GetHealthz(w http.ResponseWriter, r *http.Request) *api.Response {
	resp := api.Health{Status: "ok"}
	return api.GetHealthzJSON200Response(resp)
}

func (s Service) Hello(w http.ResponseWriter, r *http.Request) *api.Response {
	callbackURL := fmt.Sprintf("https://api.notion.com/v1/oauth/authorize?owner=user&client_id=%s&redirect_uri=%s&response_type=code", s.oauthClientID, s.redirectURI)
	_, _ = w.Write([]byte(url.QueryEscape(callbackURL)))

	return &api.Response{}
}
