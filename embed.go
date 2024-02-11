package notionstix

import "embed"

//go:embed hack/*.json
var FS embed.FS

//go:embed web/*.html
var TEMPLATES embed.FS
