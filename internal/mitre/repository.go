package mitre

import (
	"os"

	"github.com/charmbracelet/log"

	"github.com/TcM1911/stix2"
	notionstix "github.com/brittonhayes/notion-stix"
)

// NewRepository creates a new instance of the MITRE repository.
// It takes in a byte slice of STIX data and optional configuration options.
func NewRepository(data []byte, options ...Option) notionstix.Repository {
	c, err := stix2.FromJSON(data)
	if err != nil {
		panic(err)
	}

	m := MITRE{
		collection: c,
		Logger:     log.New(os.Stderr),
	}

	for _, option := range options {
		option(&m)
	}

	return &m
}
