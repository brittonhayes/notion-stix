package service

import (
	"context"
	"net/http"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/brittonhayes/notion-stix/internal/cookies"
)

type HomeData struct {
	IntegrationURL string
	Authenticated  bool
}

func (s *Service) GetEvents(w http.ResponseWriter, r *http.Request) *api.Response {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	_, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: "internal server error caused by missing bot_id cookie", Code: http.StatusInternalServerError})
	}

	// go s.subscribers[botID].Listen(w)

	// go func() {
	// 	for update := range s.updates[botID] {
	// 		fmt.Fprintf(w, "data: %s \n\n", update)
	// 		w.(http.Flusher).Flush()
	// 	}
	// }()

	<-ctx.Done()
	return nil
}

func (s *Service) GetHomePage(w http.ResponseWriter, r *http.Request) *api.Response {
	_, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	s.templates.ExecuteTemplate(w, "home", HomeData{
		Authenticated:  err == nil,
		IntegrationURL: "https://api.notion.com/v1/oauth/authorize?owner=user&client_id=080c1454-5a25-43af-b5ab-06162b1955d9&redirect_uri=https%3A%2F%2Fnotion-stix.up.railway.app%2Fauth%2Fnotion%2Fcallback&response_type=code",
	})

	return nil
}
