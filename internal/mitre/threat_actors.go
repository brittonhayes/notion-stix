package mitre

import (
	"context"
	"time"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

const (
	// threatActorDatabaseTitle is the title of the ThreatActors database.
	threatActorDatabaseTitle = "MITRE ATT&CK - ThreatActors"
	// threatActorDatabaseIcon is the icon of the ThreatActors database.
	threatActorDatabaseIcon = "⚔️"
	// threatActorPageIcon is the icon of the ThreatActor page.
	threatActorPageIcon = "⚔️"
)

// ListThreatActors returns all the ThreatActors in the MITRE collection.
func (m *MITRE) ListThreatActors() []*stix2.ThreatActor {
	return m.Collection.ThreatActors()
}

// ThreatActorByID returns the ThreatActor with the specified ID.
func (m *MITRE) ThreatActorByID(id string) *stix2.ThreatActor {
	return m.Collection.ThreatActor(stix2.Identifier(id))
}

// CreateThreatActorsDatabase creates a new ThreatActors database in Notion.
func (m *MITRE) CreateThreatActorsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: threatActorDatabaseTitle}}},
		Description:  []notion.RichText{{Text: &notion.Text{Content: "A database of MITRE ATT&CK ThreatActors of compromise."}}},
		Properties: notion.DatabaseProperties{
			"Name": {
				Type:  notion.DBPropTypeTitle,
				Title: &notion.EmptyMetadata{},
			},
			"Description": {
				Type:     notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			},
			"Sophistication": {
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
		},
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(threatActorDatabaseIcon),
		},
	}

	m.Logger.Info("Creating Notion database", "title", threatActorDatabaseTitle)
	return client.CreateDatabase(ctx, params)
}

// CreateThreatActorPage creates a new ThreatActor page in the specified ThreatActors database.
// https://oasis-open.github.io/cti-documentation/examples/defining-campaign-ta-is.html
func (m *MITRE) CreateThreatActorPage(ctx context.Context, client *notion.Client, db notion.Database, threatActor *stix2.ThreatActor) (notion.Page, error) {
	var blocks []notion.Block

	properties := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   db.ID,
		Children:   blocks,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(threatActorDatabaseIcon),
		},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: threatActor.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: threatActor.Description}},
				},
			},
			"Sophistication": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: threatActor.Sophistication}},
				},
			},
			"Motivation": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: threatActor.PrimaryMotivation}},
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(threatActor.Created.Time, false),
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
	m.Logger.Debug("Creating page", "name", threatActor.Name, "type", "ThreatActor")
	return client.CreatePage(ctx, properties)
}

type createThreatActorPageParams struct {
	ThreatActor *stix2.ThreatActor
	ParentID    string
}

func marshalThreatActor(params createThreatActorPageParams) notion.CreatePageParams {
	var blocks []notion.Block

	properties := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   params.ParentID,
		Children:   blocks,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(threatActorDatabaseIcon),
		},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: params.ThreatActor.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: params.ThreatActor.Description}},
				},
			},
			"Sophistication": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: params.ThreatActor.Sophistication}},
				},
			},
			"Motivation": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: params.ThreatActor.PrimaryMotivation}},
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

	return properties
}
