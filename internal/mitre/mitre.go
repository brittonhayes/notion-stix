// Package mitre provides functionality for working with MITRE ATT&CK data in Notion.
package mitre

import (
	"log/slog"

	"github.com/TcM1911/stix2"
	"github.com/dstotijn/go-notion"
)

// MITRE represents the API for
// integrating the MITRE ATT&CK framework in Notion.
type MITRE struct {
	collection *stix2.Collection
	Logger     *slog.Logger
}

// Option is a functional option for configuring the MITRE struct.
type Option func(*MITRE)

// WithLogger sets the logger for the MITRE struct.
func WithLogger(logger *slog.Logger) Option {
	return func(m *MITRE) {
		m.Logger = logger
	}
}

// WithCollection sets the STIX2 collection for the MITRE struct.
func WithCollection(collection *stix2.Collection) Option {
	return func(m *MITRE) {
		m.collection = collection
	}
}

// limitString truncates a string to a specified limit.
func limitString(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	return s[:limit]
}

// referencesToBlocks converts a slice of STIX2 external references to Notion blocks.
func referencesToBlocks(references []*stix2.ExternalReference) []notion.Block {
	var blocks []notion.Block

	for _, ref := range references {
		if ref.URL == "" {
			continue
		}
		blocks = append(blocks, notion.BookmarkBlock{
			URL: ref.URL,
		})
	}

	return blocks
}

func capabilitiesToBlocks(capabilities []string) []notion.Block {
	var blocks []notion.Block

	for _, capability := range capabilities {
		if capability == "" {
			continue
		}
		blocks = append(blocks, notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{Text: &notion.Text{Content: capability}},
			},
		})
	}

	return blocks
}
