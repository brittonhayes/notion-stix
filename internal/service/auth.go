package service

type OAuthAccessToken struct {
	AccessToken   string `json:"access_token,omitempty"`
	WorkspaceID   string `json:"workspace_id,omitempty"`
	WorkspaceName string `json:"workspace_name,omitempty"`
	WorkspaceIcon string `json:"workspace_icon,omitempty"`
	BotID         string `json:"bot_id,omitempty"`
}

type OAuthGrant struct {
	GrantType   string `json:"grant_type,omitempty"`
	Code        string `json:"code,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
}
