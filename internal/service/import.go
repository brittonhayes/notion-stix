package service

import (
	"context"
	"time"

	"github.com/brittonhayes/notion-stix/internal/mitre"
	"github.com/charmbracelet/log"
	"github.com/dstotijn/go-notion"
)

func (s Service) importAttackPatternsIntelToNotionDB(ctx context.Context, client *notion.Client, pageID string) error {
	limiter := time.Tick(600 * time.Millisecond)

	attackPatterns := s.repo.AttackPatterns()
	log.Info("Found attack patterns intel", "records", len(attackPatterns))

	log.Info("Creating Notion database", "title", mitre.ATTACK_PATTERN_DATABASE_TITLE)
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

func (s Service) importCampaignsIntelToNotionDB(ctx context.Context, client *notion.Client, pageID string) error {
	limiter := time.Tick(600 * time.Millisecond)

	campaigns := s.repo.Campaigns()
	log.Info("Found campaigns intel", "records", len(campaigns))

	log.Info("Creating notion database", "title", mitre.CAMPAIGNS_DATABASE_TITLE)
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

func (s Service) importMalwareIntelToNotionDB(ctx context.Context, client *notion.Client, pageID string) error {
	limiter := time.Tick(600 * time.Millisecond)

	malware := s.repo.Malware()
	log.Info("Found malware intel", "records", len(malware))

	log.Info("Creating notion database", "title", mitre.MALWARE_DATABASE_TITLE)
	malwareDB, err := s.repo.CreateMalwareDatabase(ctx, client, pageID)
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

func (s Service) importSTIXToNotion(client *notion.Client) error {

	log.Info("Creating notion pages (this might take a while)")

	ctx := context.Background()
	parentPage := "257d3f4e70f246cbad438971f691ed2d"
	err := s.importAttackPatternsIntelToNotionDB(ctx, client, parentPage)
	if err != nil {
		return err
	}

	err = s.importCampaignsIntelToNotionDB(ctx, client, parentPage)
	if err != nil {
		return err
	}

	err = s.importMalwareIntelToNotionDB(ctx, client, parentPage)
	if err != nil {
		return err
	}

	return nil
}
