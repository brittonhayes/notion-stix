package mitre

import (
	"context"
	"time"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

const (
	// campaignsDatabaseTitle is the title of the campaigns database.
	campaignsDatabaseTitle = "MITRE ATT&CK - Campaigns"
	// campaignsDatabaseIcon is the icon of the campaigns database.
	campaignsDatabaseIcon = "🗺️"
	// campaignsPageIcon is the icon of the campaign page.
	campaignsPageIcon = "🗺️"
)

// ListCampaigns returns all campaigns in the collection.
func (m *MITRE) ListCampaigns() []*stix2.Campaign {
	return m.Collection.Campaigns()
}

// CreateCampaignsDatabase creates a campaigns database in Notion.
func (m *MITRE) CreateCampaignsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: campaignsDatabaseTitle}}},
		Description:  []notion.RichText{{Text: &notion.Text{Content: "A database of MITRE ATT&CK campaigns."}}},
		Properties: notion.DatabaseProperties{
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
			"Last Seen": {
				Type: notion.DBPropTypeDate,
				Date: &notion.EmptyMetadata{},
			},
			"Created": {
				Type: notion.DBPropTypeDate,
				Date: &notion.EmptyMetadata{},
			},
			"Imported": {
				Type: notion.DBPropTypeDate,
				Date: &notion.EmptyMetadata{},
			},
		},
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(campaignsDatabaseIcon),
		},
	}

	m.Logger.Info("Creating Notion database", "title", campaignsDatabaseTitle)
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
			Emoji: notion.StringPtr(campaignsPageIcon),
		},
		Title: []notion.RichText{
			{Text: &notion.Text{Content: campaign.Name}},
		},
		Children: blocks,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: campaign.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: limitString(campaign.Description, 2000)}},
				},
			},
			"Objective": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: limitString(campaign.Objective, 2000)}},
				},
			},
			"First Seen": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(campaign.FirstSeen.Time, false),
				},
			},
			"Last Seen": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(campaign.LastSeen.Time, false),
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(campaign.Created.Time, false),
				},
			},
			"Imported": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(time.Now(), false),
				},
			},
		},
	}

	m.Logger.Debug("Creating page", "name", campaign.Name, "type", "campaign")
	return client.CreatePage(ctx, properties)
}
