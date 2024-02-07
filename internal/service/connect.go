package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brittonhayes/notion-stix/internal/api"
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
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrOAuthGrant, Code: 500})
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, NOTION_OAUTH_URL, bytes.NewReader(b))
	if err != nil {
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}
	req.SetBasicAuth(s.oauthClientID, s.oauthClientSecret)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	s.logger.Info("Requesting token from Notion API")

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error(fmt.Errorf("Notion API returned status code %d", resp.StatusCode))
		return api.ConnectJSON500Response(api.Error{Message: ErrTokenRequest, Code: 500})
	}

	var body OAuthAccessToken
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		s.logger.Error(err)
		return nil
	}

	token := body.AccessToken

	if token == "" {
		s.logger.Error("No token received from Notion API")
		return api.ConnectJSON500Response(api.Error{Message: ErrMissingToken, Code: 500})
	}

	s.logger.Info("Token received from Notion API")
	http.Redirect(w, r, "https://www.notion.so", http.StatusFound)

	client := notion.NewClient(token, notion.WithHTTPClient(s.client))
	err = s.importSTIXToNotion(client)
	if err != nil {
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: ErrImportSTIX, Code: 500})
	}

	return nil
}
