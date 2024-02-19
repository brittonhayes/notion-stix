package mitre

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// Constants for the attack pattern database and page icons.
const (
	attackPatternDatabaseTitle = "MITRE ATT&CK - Attack Patterns"
	attackPatternDatabaseIcon  = "📔"
	attackPatternPageIcon      = "📔"
)

// ATTACK_PATTERN_PROPERTIES defines the properties of the attack pattern database.
// ListAttackPatterns returns all attack patterns in the collection.
func (m *MITRE) ListAttackPatterns(collection *stix2.Collection) []*stix2.AttackPattern {
	return collection.AttackPatterns()
}

// CreateAttackPatternsDatabase creates a new attack patterns database in Notion.
func (m *MITRE) CreateAttackPatternsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: attackPatternDatabaseTitle}}},
		Description:  []notion.RichText{{Text: &notion.Text{Content: "A database of MITRE ATT&CK attack patterns."}}},
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
			},
			"Killchain Phase": {
				Type: notion.DBPropTypeMultiSelect,
				MultiSelect: &notion.SelectMetadata{
					Options: []notion.SelectOptions{
						{Name: "reconnaissance"},
						{Name: "weaponization"},
						{Name: "delivery"},
						{Name: "exploitation"},
						{Name: "installation"},
						{Name: "command-and-control"},
						{Name: "actions-on-objectives"},
					},
				},
			},
		},
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(attackPatternDatabaseIcon),
		},
	}

	// m.Logger.Info("Creating Notion database", "title", ATTACK_PATTERN_DATABASE_TITLE)
	return client.CreateDatabase(ctx, params)
}

// CreateAttackPatternPage creates a new attack pattern page in the specified database.
func (m *MITRE) CreateAttackPatternPage(ctx context.Context, client *notion.Client, databaseID string, attackPattern *stix2.AttackPattern) (notion.Page, error) {
	var blocks []notion.Block
	blocks = append(blocks, []notion.Block{
		notion.Heading2Block{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: "References"}}},
		},
	}...)

	blocks = append(blocks, referencesToBlocks(attackPattern.ExternalReferences)...)

	properties := notion.CreatePageParams{
		ParentID:   databaseID,
		ParentType: notion.ParentTypeDatabase,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(attackPatternPageIcon),
		},
		Title: []notion.RichText{
			{Text: &notion.Text{Content: attackPattern.Name}},
		},
		Children: blocks,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: attackPattern.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: limitString(attackPattern.Description, 2000)}},
				},
			},
			"Killchain Phase": notion.DatabasePageProperty{
				Type:        notion.DBPropTypeMultiSelect,
				MultiSelect: killchainPhaseToSelect(attackPattern.KillChainPhase),
			},
			"Created": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(attackPattern.Created.Time, false),
				},
			},
		},
	}

	// m.Logger.Debug("Creating page", "name", attackPattern.Name, "type", "attack-pattern")
	return client.CreatePage(ctx, properties)
}
