package service

type OAuthAccessToken struct {
	AccessToken          string `json:"access_token,omitempty"`
	WorkspaceID          string `json:"workspace_id,omitempty"`
	WorkspaceName        string `json:"workspace_name,omitempty"`
	WorkspaceIcon        string `json:"workspace_icon,omitempty"`
	BotID                string `json:"bot_id,omitempty"`
	DuplicatedTemplateID string `json:"duplicated_template_id,omitempty"`
}

type OAuthGrant struct {
	GrantType   string `json:"grant_type,omitempty"`
	Code        string `json:"code,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// TODO store access_token as value and bot_id as key in
// https://developers.notion.com/docs/authorization#step-5-the-integration-stores-the-access_token-for-future-requests
//
// I'm a bit stuck here. I need to store the access_token and then re-use that for future requests. I could put the access tokens into
// a kv store with short TTL. I could also put the access tokens into a database.
//
// The challenge is when someone goes to the website and clicks the connect to notion button, it starts an oauth flow that will timeout
// if I perform the import during that process
