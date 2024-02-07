package mitre

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// INDICATOR_DATABASE_TITLE is the title of the indicators database.
const INDICATOR_DATABASE_TITLE = "Indicators"

// INDICATOR_DATABASE_ICON is the icon of the indicators database.
const INDICATOR_DATABASE_ICON = "🔍"

// INDICATOR_PAGE_ICON is the icon of the indicator page.
const INDICATOR_PAGE_ICON = "🔍"

// INDICATOR_PROPERTIES defines the properties of the indicators database.
var INDICATOR_PROPERTIES = notion.DatabaseProperties{
	"Name": {
		Type:  notion.DBPropTypeTitle,
		Title: &notion.EmptyMetadata{},
	},
	"Description": {
		Type:     notion.DBPropTypeRichText,
		RichText: &notion.EmptyMetadata{},
	},
	"Created": {
		Type: notion.DBPropTypeDate,
		Date: &notion.EmptyMetadata{},
	},
}

// ListIndicators returns all the indicators in the MITRE collection.
func (m *MITRE) ListIndicators() []*stix2.Indicator {
	return m.collection.Indicators()
}

// indicatorByID returns the indicator with the specified ID.
func (m *MITRE) indicatorByID(id string) *stix2.Indicator {
	return m.collection.Indicator(stix2.Identifier(id))
}

// CreateIndicatorsDatabase creates a new indicators database in Notion.
func (m *MITRE) CreateIndicatorsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: INDICATOR_DATABASE_TITLE}}},
		Properties:   INDICATOR_PROPERTIES,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(INDICATOR_DATABASE_ICON),
		},
	}

	m.Logger.Info("Creating Notion database", "title", INDICATOR_DATABASE_TITLE)
	return client.CreateDatabase(ctx, params)
}

// CreateIndicatorPage creates a new indicator page in the specified indicators database.
func (m *MITRE) CreateIndicatorPage(ctx context.Context, client *notion.Client, db notion.Database, indicator *stix2.Indicator) (notion.Page, error) {
	var blocks []notion.Block

	blocks = append(blocks, []notion.Block{
		notion.Heading2Block{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: "Pattern"}}},
		},
	}...)

	blocks = append(blocks, []notion.Block{
		notion.CodeBlock{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: indicator.Pattern}}},
		},
	}...)

	properties := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   db.ID,
		Children:   blocks,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(INDICATOR_PAGE_ICON),
		},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: INDICATOR_PROPERTIES["Name"].Type,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: indicator.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: INDICATOR_PROPERTIES["Description"].Type,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: indicator.Description}},
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: INDICATOR_PROPERTIES["Created"].Type,
				Date: &notion.Date{
					Start: notion.NewDateTime(indicator.Created.Time, false),
				},
			},
		},
	}
	m.Logger.Debug("Creating page", "name", indicator.Name, "type", "indicator")
	return client.CreatePage(ctx, properties)
}
