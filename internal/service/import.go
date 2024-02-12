package service

import (
	"context"
	"net/http"
	"time"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/brittonhayes/notion-stix/internal/cookies"
	"github.com/brittonhayes/notion-stix/internal/tasks"
	"github.com/dstotijn/go-notion"
)

const (
	// FIXME this is a temporary limit to prevent the server from timing out
	// This should be removed once the task queue is implemented
	MAX_PAGES = 50
)

type authenticationResponse struct {
	pageID string
	client *notion.Client
}

func (s *Service) authenticate(w http.ResponseWriter, r *http.Request) (error, *authenticationResponse) {
	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return err, nil
	}

	pageID, err := cookies.ReadEncrypted(r, "page_id", []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return err, nil
	}

	token, err := s.store.Get(botID)
	if err != nil {
		s.logger.Error(err)
		return err, nil
	}

	client := notion.NewClient(token, notion.WithHTTPClient(s.client))

	return nil, &authenticationResponse{
		pageID: pageID,
		client: client,
	}
}

func (s *Service) ImportSTIX(w http.ResponseWriter, r *http.Request) *api.Response {
	// TODO this takes an insane amount of time. Need to implement a task queue or something.
	// Potentially also offer different import options for a subset of MITRE ATT&CK
	// Maybe use this with redis https://github.com/hibiken/asynq
	// Also maybe worth considering SSE for the client to listen for updates
	err := s.importAttackPatternsIntelToNotionDB(w, r)
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	err = s.importCampaignsIntelToNotionDB(w, r)
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	err = s.importMalwareIntelToNotionDB(w, r)
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	http.Redirect(w, r, NOTION_URL, http.StatusFound)
	return nil
}

func (s *Service) importAttackPatternsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	limiter := time.NewTicker(600 * time.Millisecond)

	err, auth := s.authenticate(w, r)
	if err != nil {
		return err
	}

	ctx := context.Background()

	attackPatternDB, err := s.repo.CreateAttackPatternsDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	attackPatterns := s.repo.ListAttackPatterns(s.repo.ListCollection())
	for i, attackPattern := range attackPatterns {
		if i > MAX_PAGES {
			return nil
		}
		<-limiter.C
		task, err := tasks.NewCreateAttackPatternsPageTask(ctx, tasks.CreateAttackPatternPagePayload{
			ParentPageID:  attackPatternDB.ID,
			AttackPattern: attackPattern,
			NotionClient:  auth.client,
		})
		if err != nil {
			s.logger.Error(err)
			return err
		}

		info, err := s.queue.Client.Enqueue(task)
		if err != nil {
			s.logger.Error(err, "failed to enqueue task", "task", task.Type)
			return err
		}
		s.logger.Info("enqueued task", "task", info.ID, "queue", info.Queue)

		// _, err = s.repo.CreateAttackPatternPage(ctx, client, attackPatternDB.ID, attackPattern)
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func (s *Service) importCampaignsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	limiter := time.NewTicker(600 * time.Millisecond)
	ctx := context.Background()

	err, auth := s.authenticate(w, r)
	if err != nil {
		return err
	}

	campaigns := s.repo.ListCampaigns()
	campaignDB, err := s.repo.CreateCampaignsDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	for i, c := range campaigns {
		if i > MAX_PAGES {
			return nil
		}
		<-limiter.C
		_, err := s.repo.CreateCampaignPage(ctx, auth.client, campaignDB, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) importMalwareIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	limiter := time.NewTicker(600 * time.Millisecond)
	ctx := context.Background()

	err, auth := s.authenticate(w, r)
	if err != nil {
		return err
	}

	malware := s.repo.ListMalware()
	malwareDB, err := s.repo.CreateMalwareDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	for i, mw := range malware {
		if i > 50 {
			return nil
		}
		<-limiter.C
		_, err = s.repo.CreateMalwarePage(ctx, auth.client, malwareDB, mw)
		if err != nil {
			return err
		}
	}
	return nil
}
