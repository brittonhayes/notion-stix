package mitre

import (
	"context"
	"time"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

const (
	// intrusionSetDatabaseTitle is the title of the IntrusionSets database.
	intrusionSetDatabaseTitle = "MITRE ATT&CK - Intrusion Sets"
	intrusionSetDatabaseIcon  = "📁"
	intrusionSetPageIcon      = "📁"
)

// ListIntrusionSets returns all the IntrusionSets in the MITRE collection.
func (m *MITRE) ListIntrusionSets(collection *stix2.Collection) []*stix2.IntrusionSet {
	return m.Collection.IntrusionSets()
}

// CreateIntrusionSetsDatabase creates a new IntrusionSets database in Notion.
func (m *MITRE) CreateIntrusionSetsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: intrusionSetDatabaseTitle}}},
		Description:  []notion.RichText{{Text: &notion.Text{Content: "A database of MITRE ATT&CK IntrusionSets."}}},
		Properties: notion.DatabaseProperties{
			"Name": {
				Type:  notion.DBPropTypeTitle,
				Title: &notion.EmptyMetadata{},
			},
			"Description": {
				Type:     notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			},
			"Motivation": {
				Type:     notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
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
			Emoji: notion.StringPtr(intrusionSetDatabaseIcon),
		},
	}

	m.Logger.Info("Creating Notion database", "title", intrusionSetDatabaseTitle)
	return client.CreateDatabase(ctx, params)
}

// CreateIntrusionSetPage creates a new IntrusionSet page in the specified IntrusionSets database.
func (m *MITRE) CreateIntrusionSetPage(ctx context.Context, client *notion.Client, databaseID string, IntrusionSet *stix2.IntrusionSet) (notion.Page, error) {
	is := intrusionSet{
		IntrusionSet: IntrusionSet,
	}

	properties := is.toNotionPageParams(databaseID)
	m.Logger.Debug("Creating page", "name", IntrusionSet.Name, "type", "IntrusionSet")
	return client.CreatePage(ctx, properties)
}

type intrusionSet struct {
	*stix2.IntrusionSet
}

func (i *intrusionSet) toNotionPageParams(parentID string) notion.CreatePageParams {
	var blocks []notion.Block

	blocks = append(blocks, []notion.Block{
		notion.Heading2Block{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: "References"}}},
		},
	}...)

	blocks = append(blocks, referencesToBlocks(i.ExternalReferences)...)

	properties := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   parentID,
		Children:   blocks,
		Icon:       &notion.Icon{Type: notion.IconTypeEmoji, Emoji: notion.StringPtr(intrusionSetPageIcon)},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type:  notion.DBPropTypeTitle,
				Title: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: i.Name}}},
			},
			"Description": notion.DatabasePageProperty{Type: notion.DBPropTypeRichText, RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: i.Description}}}},
			"Motivation":  notion.DatabasePageProperty{Type: notion.DBPropTypeRichText, RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: i.PrimaryMotivation}}}},
			"Created":     notion.DatabasePageProperty{Type: notion.DBPropTypeDate, Date: &notion.Date{Start: notion.NewDateTime(i.Created.Time, false)}},
			"Imported":    notion.DatabasePageProperty{Type: notion.DBPropTypeDate, Date: &notion.Date{Start: notion.NewDateTime(time.Now(), false)}},
		},
	}

	return properties
}
