package mitre

import (
	"context"
	"time"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

const (
	// groupDatabaseTitle is the title of the groups database.
	groupDatabaseTitle = "MITRE ATT&CK - Groups"
	groupDatabaseIcon  = "📁"
	groupPageIcon      = "📁"
)

// Listgroups returns all the groups in the MITRE collection.
func (m *MITRE) ListGroups(collection *stix2.Collection) []*stix2.IntrusionSet {
	return m.Collection.IntrusionSets()
}

// CreategroupsDatabase creates a new groups database in Notion.
func (m *MITRE) CreateGroupsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: groupDatabaseTitle}}},
		Description:  []notion.RichText{{Text: &notion.Text{Content: "A database of MITRE ATT&CK groups."}}},
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
			Emoji: notion.StringPtr(groupDatabaseIcon),
		},
	}

	m.Logger.Info("Creating Notion database", "title", groupDatabaseTitle)
	return client.CreateDatabase(ctx, params)
}

// CreateGroupPage creates a new group page in the specified groups database.
func (m *MITRE) CreateGroupPage(ctx context.Context, client *notion.Client, databaseID string, group *stix2.IntrusionSet) (notion.Page, error) {
	var blocks []notion.Block

	blocks = append(blocks, []notion.Block{
		notion.Heading2Block{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: "References"}}},
		},
	}...)

	blocks = append(blocks, referencesToBlocks(group.ExternalReferences)...)

	properties := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   databaseID,
		Children:   blocks,
		Icon:       &notion.Icon{Type: notion.IconTypeEmoji, Emoji: notion.StringPtr(groupPageIcon)},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type:  notion.DBPropTypeTitle,
				Title: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: group.Name}}},
			},
			"Description": notion.DatabasePageProperty{Type: notion.DBPropTypeRichText, RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: group.Description}}}},
			"Motivation":  notion.DatabasePageProperty{Type: notion.DBPropTypeRichText, RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: group.PrimaryMotivation}}}},
			"Created":     notion.DatabasePageProperty{Type: notion.DBPropTypeDate, Date: &notion.Date{Start: notion.NewDateTime(group.Created.Time, false)}},
			"Imported":    notion.DatabasePageProperty{Type: notion.DBPropTypeDate, Date: &notion.Date{Start: notion.NewDateTime(time.Now(), false)}},
		},
	}

	m.Logger.Debug("Creating page", "name", group.Name, "type", "group")
	return client.CreatePage(ctx, properties)
}
