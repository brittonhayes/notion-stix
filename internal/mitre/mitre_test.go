package mitre

import (
	"testing"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
	"github.com/stretchr/testify/assert"
)

func TestLimitString(t *testing.T) {
	type test struct {
		input string
		want  string
		limit int
	}

	tests := []test{
		{input: "testing", limit: 4, want: "test"},
		{input: "testing", limit: 7, want: "testing"},
		{input: "testing", limit: 0, want: ""},
	}

	for _, tc := range tests {
		got := limitString(tc.input, tc.limit)
		assert.Equal(t, got, tc.want)
	}
}

func TestReferencesToBlocks(t *testing.T) {
	type test struct {
		input []*stix2.ExternalReference
		want  []notion.Block
	}

	tests := []test{
		{
			input: []*stix2.ExternalReference{{Name: "Notion", Description: "testing", URL: "https://example.com/123"}},
			want: []notion.Block{
				notion.BookmarkBlock{URL: "https://example.com/123"},
			},
		},
		{
			input: []*stix2.ExternalReference{{Name: "Notion", Description: "testing", URL: ""}},
			want:  []notion.Block(nil),
		},
	}

	for _, tc := range tests {
		got := referencesToBlocks(tc.input)
		assert.Equal(t, got, tc.want)
	}
}

func TestCapabilitiesToBlocks(t *testing.T) {
	type test struct {
		input []string
		want  []notion.Block
	}
	tests := []test{
		{
			input: []string{"exfil", "recon"},
			want: []notion.Block{
				notion.BulletedListItemBlock{
					RichText: []notion.RichText{
						{Text: &notion.Text{Content: "exfil"}},
					},
				},
				notion.BulletedListItemBlock{
					RichText: []notion.RichText{
						{Text: &notion.Text{Content: "recon"}},
					},
				},
			},
		},
		{
			input: []string{"recon"},
			want: []notion.Block{
				notion.BulletedListItemBlock{
					RichText: []notion.RichText{
						{Text: &notion.Text{Content: "recon"}},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		got := capabilitiesToBlocks(tc.input)
		assert.Equal(t, got, tc.want)
	}
}
