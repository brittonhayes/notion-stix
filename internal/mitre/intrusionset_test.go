package mitre

import (
	"testing"
	"time"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
	"github.com/stretchr/testify/assert"
)

func TestMarshalIntrusionSet(t *testing.T) {
	arg := intrusionSet{
		IntrusionSet: &stix2.IntrusionSet{
			Name:              "APT1",
			Description:       "APT1 is a Chinese threat actor group that has been attributed to the Chinese People's Liberation Army (PLA) Third Department 12th Bureau.",
			PrimaryMotivation: "Espionage",
		},
	}
	arg.Created = &stix2.Timestamp{Time: time.Now()}

	want := &notion.CreatePageParams{
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.Name}},
				},
			},
			"Description": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.Description}},
				},
			},
			"Motivation": notion.DatabasePageProperty{
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.PrimaryMotivation}},
				},
			},
			"Created": notion.DatabasePageProperty{
				Type: notion.DBPropTypeDate,
				Date: &notion.Date{
					Start: notion.NewDateTime(arg.Created.Time, false),
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

	got := arg.toNotionPageParams("123")

	t.Run("Intrusion set page parent type is database", func(t *testing.T) {
		assert.Equal(t, got.ParentType, notion.ParentTypeDatabase)
	})

	t.Run("Intrusion set database page has wanted properties", func(t *testing.T) {
		assert.Equal(t, got.DatabasePageProperties, want.DatabasePageProperties)
	})
}