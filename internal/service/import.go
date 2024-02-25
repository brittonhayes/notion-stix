package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/brittonhayes/notion-stix/internal/cookies"
	"github.com/dstotijn/go-notion"
)

type authenticationResponse struct {
	client *notion.Client
	pageID string
}

func (s *Service) authenticate(w http.ResponseWriter, r *http.Request) (*authenticationResponse, error) {
	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	pageID, err := cookies.ReadEncrypted(r, "page_id", []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	token, err := s.store.Get(botID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	client := notion.NewClient(token, notion.WithHTTPClient(s.client))

	return &authenticationResponse{
		pageID: pageID,
		client: client,
	}, nil
}

func (s *Service) ImportSTIX(w http.ResponseWriter, r *http.Request) *api.Response {
	botID, err := cookies.Read(r, "bot_id")
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: "internal server error caused by missing bot_id cookie", Code: http.StatusInternalServerError})
	}

	err = s.importCampaignsIntelToNotionDB(w, r)
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	err = s.importGroupsIntelToNotionDB(w, r)
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	err = s.importAttackPatternsIntelToNotionDB(w, r)
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	err = s.importMalwareIntelToNotionDB(w, r)
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}

	go func() {
		s.updates[botID] <- "All records imported."
	}()

	return nil
}

func (s *Service) importAttackPatternsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	auth, err := s.authenticate(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.Read(r, "bot_id")
	if err != nil {
		return err
	}

	attackPatternDB, err := s.repo.CreateAttackPatternsDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	attackPatterns := s.repo.ListAttackPatterns(s.repo.ListCollection())
	for i, attackPattern := range attackPatterns {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err = s.repo.CreateAttackPatternPage(ctx, auth.client, attackPatternDB.ID, attackPattern)
		if err != nil {
			return err
		}

		go func() {
			s.updates[botID] <- fmt.Sprintf("Imported %d of %d attack pattern intel records", i, len(attackPatterns))
		}()

		if i%10 == 0 || i == len(attackPatterns)-1 {
			s.logger.Info("imported attack patterns intel", "done", i, "total", len(attackPatterns))
		}
	}

	return nil
}

func (s *Service) importGroupsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	auth, err := s.authenticate(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.Read(r, "bot_id")
	if err != nil {
		return err
	}

	db, err := s.repo.CreateGroupsDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	groups := s.repo.ListGroups(s.repo.ListCollection())
	for i, group := range groups {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err = s.repo.CreateGroupPage(ctx, auth.client, db.ID, group)
		if err != nil {
			return err
		}

		go func() {
			s.updates[botID] <- fmt.Sprintf("Imported %d of %d APT intel records", i, len(groups))
		}()

		if i%10 == 0 || i == len(groups)-1 {
			s.logger.Info("imported groups intel", "done", i, "total", len(groups))
		}
	}

	return nil
}

func (s *Service) importCampaignsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	auth, err := s.authenticate(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.Read(r, "bot_id")
	if err != nil {
		return err
	}

	campaigns := s.repo.ListCampaigns()
	campaignDB, err := s.repo.CreateCampaignsDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	for i, campaign := range campaigns {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err := s.repo.CreateCampaignPage(ctx, auth.client, campaignDB, campaign)
		if err != nil {
			return err
		}

		go func() {
			s.updates[botID] <- fmt.Sprintf("Imported %d of %d campaign intel records", i, len(campaigns))
		}()

		if i%10 == 0 || i == len(campaigns)-1 {
			s.logger.Info("imported campaign intel", "done", i, "total", len(campaigns))
		}
	}
	return nil
}

func (s *Service) importMalwareIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	auth, err := s.authenticate(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.Read(r, "bot_id")
	if err != nil {
		return err
	}

	malware := s.repo.ListMalware()
	malwareDB, err := s.repo.CreateMalwareDatabase(ctx, auth.client, auth.pageID)
	if err != nil {
		return err
	}

	for i, mw := range malware {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err = s.repo.CreateMalwarePage(ctx, auth.client, malwareDB, mw)
		if err != nil {
			return err
		}

		go func() {
			s.updates[botID] <- fmt.Sprintf("Imported %d of %d malware intel records", i, len(malware))
		}()

		if i%10 == 0 || i == len(malware)-1 {
			s.logger.Info("imported malware intel", "done", i, "total", len(malware))
		}
	}
	return nil
}
