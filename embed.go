package notionstix

import (
	"embed"
)

//go:embed hack/*.json
var FS embed.FS

type StixSource int

const (
	MitreEnterpriseAttack StixSource = iota + 1
)

func (s StixSource) String() string {
	return [...]string{"./hack/enterprise-attack-14.1.json"}[s-1]
}
