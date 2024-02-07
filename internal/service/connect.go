package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/charmbracelet/log"
	"github.com/dstotijn/go-notion"
)

// Connect handles the integration setup request.
// It performs the OAuth token exchange with Notion API and returns an API response.
func (s Service) Connect(w http.ResponseWriter, r *http.Request, params api.ConnectParams) *api.Response {
	b, err := json.Marshal(&OAuthGrant{
		GrantType:   "authorization_code",
		Code:        params.Code,
		RedirectURI: s.redirectURI,
	})
	if err != nil {
		log.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrOAuthGrant, Code: 500})
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, NOTION_OAUTH_URL, bytes.NewReader(b))
	if err != nil {
		log.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}
	req.SetBasicAuth(s.oauthClientID, s.oauthClientSecret)
	req.Header.Add("Content-Type", "application/json")

	log.Info("Requesting token from Notion API")
	rsp, err := s.client.Do(req)
	if err != nil {
		log.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}
	defer rsp.Body.Close()

	var oauthResponse OAuthAccessToken
	if err = json.NewDecoder(rsp.Body).Decode(&oauthResponse); err != nil {
		log.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrTokenDecode, Code: 500})
	}
	log.Info("Token received from Notion API")

	client := notion.NewClient(oauthResponse.AccessToken, notion.WithHTTPClient(s.client))
	err = s.importSTIXToNotion(client)
	if err != nil {
		log.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrImportSTIX, Code: 500})
	}

	return nil
}
