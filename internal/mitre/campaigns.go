package mitre

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// CAMPAIGNS_DATABASE_TITLE is the title of the campaigns database.
const CAMPAIGNS_DATABASE_TITLE = "Campaigns"

// CAMPAIGNS_DATABASE_ICON is the icon of the campaigns database.
const CAMPAIGNS_DATABASE_ICON = "🗺️"

// CAMPAIGNS_PAGE_ICON is the icon of the campaign page.
const CAMPAIGNS_PAGE_ICON = "🗺️"

// CAMPAIGN_PROPERTIES defines the properties of a campaign in the Notion database.
var CAMPAIGN_PROPERTIES = notion.DatabaseProperties{
	"Name": {
		Type:  notion.DBPropTypeTitle,
		Title: &notion.EmptyMetadata{},
	},
	"Description": {
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"Objective": {
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"First Seen": {
		Type: notion.DBPropTypeDate,
		Date: &notion.EmptyMetadata{},
	},
	"Created": {
		Type: notion.DBPropTypeDate,
		Date: &notion.EmptyMetadata{},
	},
}

// ListCampaigns returns all campaigns in the collection.
func (m *MITRE) ListCampaigns() []*stix2.Campaign {
	return m.collection.Campaigns()
}

// campaignByID returns the campaign with the specified ID.
func (m *MITRE) campaignByID(id string) *stix2.Campaign {
	return m.collection.Campaign(stix2.Identifier(id))
}

// CreateCampaignsDatabase creates a campaigns database in Notion.
func (m *MITRE) CreateCampaignsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: CAMPAIGNS_DATABASE_TITLE}}},
		Properties:   CAMPAIGN_PROPERTIES,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(CAMPAIGNS_DATABASE_ICON),
		},
	}

	m.Logger.Info("Creating Notion database", "title", CAMPAIGNS_DATABASE_TITLE)
	return client.CreateDatabase(ctx, params)
}

// CreateCampaignPage creates a campaign page in Notion.
func (m *MITRE) CreateCampaignPage(ctx context.Context, client *notion.Client, db notion.Database, campaign *stix2.Campaign) (notion.Page, error) {
	var blocks []notion.Block

	blocks = append(blocks, []notion.Block{
		notion.Heading2Block{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: "References"}}},
		},
	}...)

	blocks = append(blocks, referencesToBlocks(campaign.ExternalReferences)...)

	properties := notion.CreatePageParams{
		ParentID:   db.ID,
		ParentType: notion.ParentTypeDatabase,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(CAMPAIGNS_PAGE_ICON),
		},
		Title: []notion.RichText{
			{Text: &notion.Text{Content: campaign.Name}},
		},
		Children: blocks,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: CAMPAIGN_PROPERTIES["Name"].Type,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: campaign.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: CAMPAIGN_PROPERTIES["Description"].Type,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: limitString(campaign.Description, 2000)}},
				},
			},
			"Objective": notion.DatabasePageProperty{
				Type: CAMPAIGN_PROPERTIES["Objective"].Type,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: limitString(campaign.Objective, 2000)}},
				},
			},
			"First Seen": notion.DatabasePageProperty{
				Type: CAMPAIGN_PROPERTIES["First Seen"].Type,
				Date: &notion.Date{
					Start: notion.NewDateTime(campaign.FirstSeen.Time, false),
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: CAMPAIGN_PROPERTIES["Created"].Type,
				Date: &notion.Date{
					Start: notion.NewDateTime(campaign.Created.Time, false),
				},
			},
		},
	}

	m.Logger.Debug("Creating page", "name", campaign.Name, "type", "campaign")
	return client.CreatePage(ctx, properties)
}
