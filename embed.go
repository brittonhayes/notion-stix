package notionstix

import "embed"

// TODO: make embedding optional

//go:embed hack/*.json
var FS embed.FS

//go:embed web/*.html
var TEMPLATES embed.FS
