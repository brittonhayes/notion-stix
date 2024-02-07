package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/charmbracelet/log"
	"github.com/dstotijn/go-notion"
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

	redirectURL       string
	oauthClientID     string
	oauthClientSecret string
}

// New creates a new instance of the Service.
func New(repo notionstix.Repository) Service {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = retryablehttp.LinearJitterBackoff
	retryClient.Logger = nil
	return Service{repo: repo, client: retryClient.StandardClient()}
}

// Setup handles the integration setup request.
// It performs the OAuth token exchange with Notion API and returns an API response.
func (s Service) Setup(w http.ResponseWriter, r *http.Request, params api.SetupParams) *api.Response {
	b, err := json.Marshal(&OAuthGrant{
		GrantType:   "authorization_code",
		Code:        params.Code,
		RedirectURI: s.redirectURL,
	})
	if err != nil {
		log.Error(err)
		return api.SetupJSON500Response(api.Error{Message: ErrOAuthGrant, Code: 500})
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, NOTION_OAUTH_URL, bytes.NewReader(b))
	if err != nil {
		log.Error(err)
		return api.SetupJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}
	req.SetBasicAuth(s.oauthClientID, s.oauthClientSecret)
	req.Header.Add("Content-Type", "application/json")

	rsp, err := s.client.Do(req)
	if err != nil {
		log.Error(err)
		return api.SetupJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}
	defer rsp.Body.Close()

	var body OAuthAccessToken
	if err = json.NewDecoder(rsp.Body).Decode(&body); err != nil {
		log.Error(err)
		return api.SetupJSON500Response(api.Error{Message: ErrTokenDecode, Code: 500})
	}

	client := notion.NewClient(body.AccessToken, notion.WithHTTPClient(s.client))
	err = s.importSTIXToNotion(client)
	if err != nil {
		log.Error(err)
		return api.SetupJSON500Response(api.Error{Message: ErrImportSTIX, Code: 500})
	}

	return nil
}
