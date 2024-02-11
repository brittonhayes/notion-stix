package mitre

import "github.com/TcM1911/stix2"

func (m *MITRE) ListCollection() *stix2.Collection {
	return m.Collection
}
