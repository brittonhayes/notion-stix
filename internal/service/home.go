package service

import (
	"html/template"
	"net/http"

	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/brittonhayes/notion-stix/internal/cookies"
)

type HomeData struct {
	Authenticated  bool
	IntegrationURL string
}

func (s *Service) GetHomePage(w http.ResponseWriter, r *http.Request) *api.Response {
	tmpl := template.Must(template.ParseFS(notionstix.TEMPLATES, "web/*.html"))

	botCookie, _ := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	s.logger.Info("Bot Cookie", botCookie != "")

	tmpl.ExecuteTemplate(w, "home", HomeData{
		Authenticated:  botCookie != "",
		IntegrationURL: "https://api.notion.com/v1/oauth/authorize?owner=user&client_id=080c1454-5a25-43af-b5ab-06162b1955d9&redirect_uri=https%3A%2F%2Fnotion-stix.up.railway.app%2Fauth%2Fnotion%2Fcallback&response_type=code",
	})

	return nil
}
