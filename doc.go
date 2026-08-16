// Package scaffold_go_api is the composition root's one dependency: an
// http.Handler satisfying api/openapi.yaml. Everything worth importing lives
// here, at the module root, rather than under internal/ — there's nothing
// yet worth hiding, and cmd/scaffold-go-api/main.go is the only importer.
// See go.md's module layout progression for when that changes.
package scaffold_go_api
