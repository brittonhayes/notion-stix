package mitre

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

const (
	// indicatorPageIcon is the icon of the indicator page.
	indicatorPageIcon = "🔍"
	// indicatorDatabaseTitle is the title of the indicators database.
	indicatorDatabaseTitle = "Indicators"
	// indicatorDatabaseIcon is the icon of the indicators database.
	indicatorDatabaseIcon = "🔍"
)

// ListIndicators returns all the indicators in the MITRE collection.
func (m *MITRE) ListIndicators() []*stix2.Indicator {
	return m.Collection.Indicators()
}

// indicatorByID returns the indicator with the specified ID.
func (m *MITRE) indicatorByID(id string) *stix2.Indicator {
	return m.Collection.Indicator(stix2.Identifier(id))
}

// CreateIndicatorsDatabase creates a new indicators database in Notion.
func (m *MITRE) CreateIndicatorsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: indicatorDatabaseTitle}}},
		Description:  []notion.RichText{{Text: &notion.Text{Content: "A database of MITRE ATT&CK indicators of compromise."}}},
		Properties: notion.DatabaseProperties{
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
			}},
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(indicatorDatabaseIcon),
		},
	}

	m.Logger.Info("Creating Notion database", "title", indicatorDatabaseTitle)
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
			Emoji: notion.StringPtr(indicatorPageIcon),
		},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: indicator.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: indicator.Description}},
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(indicator.Created.Time, false),
				},
			},
		},
	}
	m.Logger.Debug("Creating page", "name", indicator.Name, "type", "indicator")
	return client.CreatePage(ctx, properties)
}
