package mitre

import (
	"testing"
	"time"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
	"github.com/stretchr/testify/assert"
)

func TestMarshalThreatActor(t *testing.T) {
	arg := createThreatActorPageParams{
		ParentID: "1234",
		ThreatActor: &stix2.ThreatActor{
			Name:              "APT1",
			Description:       "APT1 is a Chinese threat actor group that has been attributed to the Chinese People's Liberation Army (PLA) Third Department 12th Bureau.",
			PrimaryMotivation: "Espionage",
			Sophistication:    "Advanced",
		},
	}

	want := &notion.DatabasePageProperties{
		"Name": notion.DatabasePageProperty{
			Type: notion.DBPropTypeTitle,
			Title: []notion.RichText{
				{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.ThreatActor.Name}},
			},
		},
		"Description": notion.DatabasePageProperty{
			Type: notion.DBPropTypeRichText,
			RichText: []notion.RichText{
				{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.ThreatActor.Description}},
			},
		},
		"Sophistication": notion.DatabasePageProperty{
			Type: notion.DBPropTypeRichText,
			RichText: []notion.RichText{
				{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.ThreatActor.Sophistication}},
			},
		},
		"Motivation": notion.DatabasePageProperty{
			Type: notion.DBPropTypeRichText,
			RichText: []notion.RichText{
				{Type: notion.RichTextTypeText, Text: &notion.Text{Content: arg.ThreatActor.PrimaryMotivation}},
			},
		},
		"Imported": notion.DatabasePageProperty{
			Type: notion.DBPropTypeDate,
			Date: &notion.Date{
				Start: notion.NewDateTime(time.Now(), false),
			},
		},
	}

	page := marshalThreatActor(arg)
	assert.Equal(t, page.DatabasePageProperties, want)
}
