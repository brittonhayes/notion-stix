package mitre

import (
	"context"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// Constants for the attack pattern database and page icons.
const (
	ATTACK_PATTERN_DATABASE_TITLE = "Attack Patterns"
	ATTACK_PATTERN_DATABASE_ICON  = "📔"
	ATTACK_PATTERN_PAGE_ICON      = "📔"
)

// ATTACK_PATTERN_PROPERTIES defines the properties of the attack pattern database.
var ATTACK_PATTERN_PROPERTIES = notion.DatabaseProperties{
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

// ListAttackPatterns returns all attack patterns in the collection.
func (m *MITRE) ListAttackPatterns() []*stix2.AttackPattern {
	return m.collection.AttackPatterns()
}

// attackPatternByID returns the attack pattern with the specified ID.
func (m *MITRE) attackPatternByID(id string) *stix2.AttackPattern {
	return m.collection.AttackPattern(stix2.Identifier(id))
}

// CreateAttackPatternsDatabase creates a new attack patterns database in Notion.
func (m *MITRE) CreateAttackPatternsDatabase(ctx context.Context, client *notion.Client, parentPageID string) (notion.Database, error) {
	params := notion.CreateDatabaseParams{
		ParentPageID: parentPageID,
		Title:        []notion.RichText{{Text: &notion.Text{Content: ATTACK_PATTERN_DATABASE_TITLE}}},
		Properties:   ATTACK_PATTERN_PROPERTIES,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(ATTACK_PATTERN_DATABASE_ICON),
		},
	}

	m.Logger.Info("Creating Notion database", "title", ATTACK_PATTERN_DATABASE_TITLE)
	return client.CreateDatabase(ctx, params)
}

// CreateAttackPatternPage creates a new attack pattern page in the specified database.
func (m *MITRE) CreateAttackPatternPage(ctx context.Context, client *notion.Client, db notion.Database, attackPattern *stix2.AttackPattern) (notion.Page, error) {

	var blocks []notion.Block
	blocks = append(blocks, []notion.Block{
		notion.Heading2Block{
			RichText: []notion.RichText{{Type: notion.RichTextTypeText, Text: &notion.Text{Content: "References"}}},
		},
	}...)

	blocks = append(blocks, referencesToBlocks(attackPattern.ExternalReferences)...)

	properties := notion.CreatePageParams{
		ParentID:   db.ID,
		ParentType: notion.ParentTypeDatabase,
		Icon: &notion.Icon{
			Type:  notion.IconTypeEmoji,
			Emoji: notion.StringPtr(ATTACK_PATTERN_PAGE_ICON),
		},
		Title: []notion.RichText{
			{Text: &notion.Text{Content: attackPattern.Name}},
		},
		Children: blocks,
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: ATTACK_PATTERN_PROPERTIES["Name"].Type,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: attackPattern.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: ATTACK_PATTERN_PROPERTIES["Description"].Type,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: limitString(attackPattern.Description, 2000)}},
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: ATTACK_PATTERN_PROPERTIES["Created"].Type,
				Date: &notion.Date{
					Start: notion.NewDateTime(attackPattern.Created.Time, false),
				},
			},
		},
	}

	m.Logger.Debug("Creating page", "name", attackPattern.Name, "type", "attack-pattern")
	return client.CreatePage(ctx, properties)
}
