package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/brittonhayes/notion-stix/internal/cookies"
)

var (
	ErrValueTooLong = "cookie value too long"
	ErrInvalidValue = "invalid cookie value"
)

type connectionRecord struct {
	Token        string `json:"token"`
	ParentPageID string `json:"parent_page_id"`
}

// Connect handles the connection request from the client.
// The access token is then used to redirect the client to the Notion URL.
// If any errors occur during the process, appropriate error responses are returned.
func (s *Service) Connect(w http.ResponseWriter, r *http.Request, params api.ConnectParams) *api.Response {
	if params.Error != nil {
		s.logger.Error(params.Error)
		return api.ConnectJSON500Response(api.Error{Message: ErrCancel, Code: 500})
	}

	if params.Code == nil {
		s.logger.Error("No code received from client")
		return api.ConnectJSON500Response(api.Error{Message: ErrOAuthGrant, Code: 500})
	}

	b, err := json.Marshal(&OAuthGrant{
		GrantType:   "authorization_code",
		Code:        *params.Code,
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
		s.logger.Error(fmt.Errorf("notion api returned status code %d", resp.StatusCode))
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
		return api.ConnectJSON500Response(api.Error{Message: ErrMissingToken, Code: http.StatusBadRequest})
	}

	b, err = json.Marshal(connectionRecord{Token: token, ParentPageID: body.DuplicatedTemplateID})
	if err != nil {
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: err.Error(), Code: http.StatusBadRequest})
	}

	s.logger.Info("Token received from Notion API")
	err = s.store.Set(body.BotID, b)
	if err != nil {
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: err.Error(), Code: http.StatusBadRequest})
	}

	botCookie := http.Cookie{
		Name:     "bot_id",
		Value:    body.BotID,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	err = cookies.WriteEncrypted(w, botCookie, []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return api.ConnectJSON500Response(api.Error{Message: err.Error(), Code: http.StatusBadRequest})
	}

	go func() {
		s.updates[body.BotID] <- "Connected"
	}()

	// TODO: store the status of the import in the kv store so we can display it to the user
	http.Redirect(w, r, NOTION_URL, http.StatusFound)
	return nil
}
