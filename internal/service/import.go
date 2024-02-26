package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/brittonhayes/notion-stix/internal/api"
	"github.com/brittonhayes/notion-stix/internal/cookies"
	"github.com/brittonhayes/notion-stix/internal/kv"
	"github.com/dstotijn/go-notion"
)

type session struct {
	client *notion.Client
	pageID string
}

func (s *Service) newSession(w http.ResponseWriter, r *http.Request) (*session, error) {
	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	value, err := s.store.Get(botID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	var rec connectionRecord
	err = json.Unmarshal(value, &rec)
	if err != nil {
		return nil, err
	}

	client := notion.NewClient(rec.Token, notion.WithHTTPClient(s.client))
	return &session{
		client: client,
		pageID: rec.ParentPageID,
	}, nil
}

func (s *Service) ImportSTIX(w http.ResponseWriter, r *http.Request) *api.Response {
	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		return api.ImportSTIXJSON500Response(api.Error{Message: "internal server error caused by missing bot_id cookie", Code: http.StatusInternalServerError})
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.updates[botID] <- "Starting import..."
		wg.Done()
	}()

	err = s.importCampaignsIntelToNotionDB(w, r)
	if err != nil {
		s.logger.Error(err)
		return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	}
	//
	// err = s.importIntrusionSetsIntelToNotionDB(w, r)
	// if err != nil {
	// 	s.logger.Error(err)
	// 	return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	// }
	//
	// err = s.importAttackPatternsIntelToNotionDB(w, r)
	// if err != nil {
	// 	s.logger.Error(err)
	// 	return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	// }
	//
	// err = s.importMalwareIntelToNotionDB(w, r)
	// if err != nil {
	// 	s.logger.Error(err)
	// 	return api.ImportSTIXJSON500Response(api.Error{Message: ErrImportSTIX, Code: http.StatusInternalServerError})
	// }

	wg.Add(1)
	go func() {
		s.updates[botID] <- "All records imported."
	}()

	wg.Wait()

	return nil
}

func (s *Service) importAttackPatternsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	sess, err := s.newSession(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		return err
	}

	attackPatternDB, err := s.repo.CreateAttackPatternsDatabase(ctx, sess.client, sess.pageID)
	if err != nil {
		return err
	}

	attackPatterns := s.repo.ListAttackPatterns(s.repo.ListCollection())
	for i, attackPattern := range attackPatterns {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err := s.repo.CreateAttackPatternPage(ctx, sess.client, attackPatternDB.ID, attackPattern)
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

func (s *Service) importIntrusionSetsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	sess, err := s.newSession(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		return err
	}

	db, err := s.repo.CreateIntrusionSetsDatabase(ctx, sess.client, sess.pageID)
	if err != nil {
		return err
	}

	intrusionSets := s.repo.ListIntrusionSets(s.repo.ListCollection())
	for i, intrusionSet := range intrusionSets {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err = s.repo.CreateIntrusionSetPage(ctx, sess.client, db.ID, intrusionSet)
		if err != nil {
			return err
		}

		go func() {
			s.updates[botID] <- fmt.Sprintf("Imported %d of %d APT Intrusion Set intel records", i, len(intrusionSets))
		}()

		if i%10 == 0 || i == len(intrusionSets)-1 {
			s.logger.Info("imported IntrusionSets intel", "done", i, "total", len(intrusionSets))
		}
	}

	return nil
}

func (s *Service) importCampaignsIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	sess, err := s.newSession(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		return err
	}

	// Check if the campaign database already exists in the kv store
	// if it does, return early
	_, err = s.store.Get(fmt.Sprintf("%s-%s-%s", botID, sess.pageID, "campaigns"))
	if err == nil {
		return nil
	}

	campaignDB, err := s.repo.CreateCampaignsDatabase(ctx, sess.client, sess.pageID)
	if err != nil {
		return err
	}

	campaignDBJSON, err := json.Marshal(campaignDB)
	if err != nil {
		return err
	}

	err = s.store.Set(fmt.Sprintf("%s-%s-%s", botID, sess.pageID, "campaigns"), campaignDBJSON)
	if err != nil {
		return err
	}

	campaigns := s.repo.ListCampaigns()
	for i, campaign := range campaigns {

		_, err := s.store.Get(fmt.Sprintf("%s-%s-%s", botID, sess.pageID, campaign.ID))
		if err == kv.ErrKeyNotFound {
			// Create the campaign page and store the result in the KV store
			// if it doesn't exist
			//
			r := s.limiter.Reserve()
			time.Sleep(r.Delay())
			page, err := s.repo.CreateCampaignPage(ctx, sess.client, campaignDB, campaign)
			if err != nil {
				return err
			}

			b, err := json.Marshal(page)
			if err != nil {
				return err
			}

			err = s.store.Set(fmt.Sprintf("%s-%s-%s", botID, sess.pageID, campaign.ID), b)
			if err != nil {
				return err
			}

			go func() {
				s.logger.Info("updating bot", "id", botID, "message", fmt.Sprintf("Imported %d of %d campaign intel records", i, len(campaigns)))
				s.updates[botID] <- fmt.Sprintf("Imported %d of %d campaign intel records", i, len(campaigns))
			}()

			if i%10 == 0 || i == len(campaigns)-1 {
				s.logger.Info("imported campaign intel", "done", i, "total", len(campaigns))
			}
		}
	}
	return nil
}

func (s *Service) importMalwareIntelToNotionDB(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	sess, err := s.newSession(w, r)
	if err != nil {
		return err
	}

	botID, err := cookies.ReadEncrypted(r, "bot_id", []byte(s.cookieSecret))
	if err != nil {
		return err
	}

	malware := s.repo.ListMalware()
	malwareDB, err := s.repo.CreateMalwareDatabase(ctx, sess.client, sess.pageID)
	if err != nil {
		return err
	}

	for i, mw := range malware {
		r := s.limiter.Reserve()
		time.Sleep(r.Delay())

		_, err = s.repo.CreateMalwarePage(ctx, sess.client, malwareDB, mw)
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
