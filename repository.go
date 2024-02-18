// Package stix is the parent package to the Notion STIX integration, API, and CLI.
//
//go:generate goapi-gen -generate types,server,spec -package api --out internal/api/api.gen.go ./internal/api/openapi.yaml
package notionstix

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// Repository defines the interface for interacting with the Notion database.
type Repository interface {
	CollectionRepository
	AttackPatternRepository
	CampaignRepository
	MalwareRepository
	GroupRepository
}

// Store is the interface that defines the methods for a key-value store.
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Cleanup()
}

type AttackPatternRepository interface {
	// ListAttackPatterns returns a slice of AttackPattern objects.
	ListAttackPatterns(collection *stix2.Collection) []*stix2.AttackPattern
	// CreateAttackPatternsDatabase creates a new Notion database for AttackPatterns.
	CreateAttackPatternsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateAttackPatternPage creates a new Notion page for a specific AttackPattern.
	CreateAttackPatternPage(ctx context.Context, client *notion.Client, databaseID string, attackPattern *stix2.AttackPattern) (notion.Page, error)
}

type GroupRepository interface {
	// ListGroups returns a slice of Group objects.
	ListGroups(collection *stix2.Collection) []*stix2.IntrusionSet
	// CreateGroupsDatabase creates a new Notion database for Groups.
	CreateGroupsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateGroupPage creates a new Notion page for a specific Group.
	CreateGroupPage(ctx context.Context, client *notion.Client, databaseID string, group *stix2.IntrusionSet) (notion.Page, error)
}

type CampaignRepository interface {
	// ListCampaigns returns a slice of Campaign objects.
	ListCampaigns() []*stix2.Campaign
	// CreateCampaignsDatabase creates a new Notion database for Campaigns.
	CreateCampaignsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateCampaignPage creates a new Notion page for a specific Campaign.
	CreateCampaignPage(ctx context.Context, client *notion.Client, db notion.Database, campaign *stix2.Campaign) (notion.Page, error)
}

type IndicatorRepository interface {
	// ListIndicators returns a slice of Indicator objects.
	ListIndicators() []*stix2.Indicator
	// CreateIndicatorsDatabase creates a new Notion database for Indicators.
	CreateIndicatorsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateIndicatorPage creates a new Notion page for a specific Indicator.
	CreateIndicatorPage(ctx context.Context, client *notion.Client, db notion.Database, indicator *stix2.Indicator) (notion.Page, error)
}

type MalwareRepository interface {
	// ListMalware returns a slice of Malware objects.
	ListMalware() []*stix2.Malware
	// CreateMalwareDatabase creates a new Notion database for Malware.
	CreateMalwareDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateMalwarePage creates a new Notion page for a specific Malware.
	CreateMalwarePage(ctx context.Context, client *notion.Client, db notion.Database, malware *stix2.Malware) (notion.Page, error)
}

type CollectionRepository interface {
	// ListCollection returns the entire collection of STIX objects.
	ListCollection() *stix2.Collection
}
