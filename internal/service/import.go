package service

import (
	"context"
	"net/http"
	"time"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/dstotijn/go-notion"
)

func (s *Service) importAttackPatternsIntelToNotionDB(ctx context.Context, client *notion.Client, pageID string) error {
	limiter := time.Tick(600 * time.Millisecond)

	attackPatterns := s.repo.ListAttackPatterns()
	attackPatternDB, err := s.repo.CreateAttackPatternsDatabase(ctx, client, pageID)
	if err != nil {
		return err
	}

	for _, ap := range attackPatterns {
		<-limiter
		_, err = s.repo.CreateAttackPatternPage(ctx, client, attackPatternDB, ap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) importCampaignsIntelToNotionDB(ctx context.Context, client *notion.Client, pageID string) error {
	limiter := time.Tick(600 * time.Millisecond)

	campaigns := s.repo.ListCampaigns()
	campaignDB, err := s.repo.CreateCampaignsDatabase(ctx, client, pageID)
	if err != nil {
		return err
	}

	for _, c := range campaigns {
		<-limiter
		_, err := s.repo.CreateCampaignPage(ctx, client, campaignDB, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) importMalwareIntelToNotionDB(ctx context.Context, client *notion.Client, parentPageID string) error {
	limiter := time.Tick(600 * time.Millisecond)

	malware := s.repo.ListMalware()
	malwareDB, err := s.repo.CreateMalwareDatabase(ctx, client, parentPageID)
	if err != nil {
		return err
	}

	for _, mw := range malware {
		<-limiter
		_, err = s.repo.CreateMalwarePage(ctx, client, malwareDB, mw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) importSTIXToNotion(client *notion.Client, parentPageID string) error {

	s.logger.Info("Creating notion pages (this might take a while)")

	ctx := context.Background()
	err := s.importAttackPatternsIntelToNotionDB(ctx, client, parentPageID)
	if err != nil {
		return err
	}

	err = s.importCampaignsIntelToNotionDB(ctx, client, parentPageID)
	if err != nil {
		return err
	}

	err = s.importMalwareIntelToNotionDB(ctx, client, parentPageID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ImportSTIX(w http.ResponseWriter, r *http.Request) *api.Response {
	botID, err := r.Cookie("bot_id")
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: "missing bot_id cookie", Code: 500})
	}

	pageID, err := r.Cookie("page_id")
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: "missing page_id cookie", Code: 500})
	}

	client := notion.NewClient(s.tokens[botID.Value], notion.WithHTTPClient(s.client))

	err = s.importSTIXToNotion(client, pageID.Value)
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: "failed to import STIX data to Notion", Code: 500})
	}

	http.Redirect(w, r, NOTION_URL, http.StatusFound)

	return nil
}
