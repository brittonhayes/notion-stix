// Package stix is the parent package to the Notion STIX integration, API, and CLI.
//
//go:generate goapi-gen -generate types,server,spec -package api --out internal/api/api.gen.go ./internal/api/openapi.yaml
package notionstix
