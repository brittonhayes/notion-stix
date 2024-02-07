package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/dstotijn/go-notion"
)

// Connect handles the connection request from the client.
// It receives the authorization code from the client and exchanges it for an access token from the Notion API.
// The access token is then used to redirect the client to the Notion URL.
// If any errors occur during the process, appropriate error responses are returned.
func (s *Service) Connect(w http.ResponseWriter, r *http.Request, params api.ConnectParams) *api.Response {
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

	if token == "foo" {
		s.logger.Error("No token received from Notion API")
		u, _ := url.Parse(s.redirectURI)
		u.Query().Add("error", "missing token in request body")

		http.Redirect(w, r, u.String(), http.StatusFound)
		return api.ConnectJSON500Response(api.Error{Message: ErrMissingToken, Code: 500})
	}

	s.logger.Info("Token received from Notion API")
	s.tokens[body.BotID] = token

	s.logger.Info("Starting notion import for bot", "bot_id", body.BotID)

	//FIXME this is a temporary approach to start a goroutine for the import
	go s.startNotionImport(body.BotID, body.DuplicatedTemplateID)

	http.Redirect(w, r, NOTION_URL, http.StatusFound)

	return nil
}

func (s *Service) startNotionImport(botID string, parentPageID string) {
	// Retrieve the token from storage
	token, exists := s.tokens[botID]
	if !exists {
		s.logger.Error("Token not found for botID:", botID)
		return
	}

	// Create a new Notion client with the token
	client := notion.NewClient(token, notion.WithHTTPClient(s.client))

	// Perform the import operation
	err := s.importSTIXToNotion(client, parentPageID)
	if err != nil {
		s.logger.Error(err)
	} else {
		s.logger.Info("STIX import to Notion completed successfully")
	}
}
