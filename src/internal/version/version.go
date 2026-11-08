// Package version holds build information injected through -ldflags.
package version

// Build information, set by the Makefile and goreleaser.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
