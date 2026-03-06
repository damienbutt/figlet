// Package version holds build-time metadata injected via ldflags.
package version

// Version, GitCommit, and BuildDate are set at build time via -ldflags.
var (
	Version   = "dev"
	GitCommit = ""
	BuildDate = ""
)
