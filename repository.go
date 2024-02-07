package notionstix

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// Repository defines the interface for interacting with the Notion database.
type Repository interface {
	// ListAttackPatterns returns a slice of AttackPattern objects.
	ListAttackPatterns() []*stix2.AttackPattern
	// CreateAttackPatternsDatabase creates a new Notion database for AttackPatterns.
	CreateAttackPatternsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateAttackPatternPage creates a new Notion page for a specific AttackPattern.
	CreateAttackPatternPage(ctx context.Context, client *notion.Client, db notion.Database, attackPattern *stix2.AttackPattern) (notion.Page, error)

	// ListCampaigns returns a slice of Campaign objects.
	ListCampaigns() []*stix2.Campaign
	// CreateCampaignsDatabase creates a new Notion database for Campaigns.
	CreateCampaignsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateCampaignPage creates a new Notion page for a specific Campaign.
	CreateCampaignPage(ctx context.Context, client *notion.Client, db notion.Database, campaign *stix2.Campaign) (notion.Page, error)

	// ListIndicators returns a slice of Indicator objects.
	ListIndicators() []*stix2.Indicator
	// CreateIndicatorsDatabase creates a new Notion database for Indicators.
	CreateIndicatorsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateIndicatorPage creates a new Notion page for a specific Indicator.
	CreateIndicatorPage(ctx context.Context, client *notion.Client, db notion.Database, indicator *stix2.Indicator) (notion.Page, error)

	// ListMalware returns a slice of ListMalware objects.
	ListMalware() []*stix2.Malware
	// CreateMalwareDatabase creates a new Notion database for Malware.
	CreateMalwareDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error)
	// CreateMalwarePage creates a new Notion page for a specific Malware.
	CreateMalwarePage(ctx context.Context, client *notion.Client, db notion.Database, malware *stix2.Malware) (notion.Page, error)
}
