package figlet

import "embed"

// FontFS is the embedded filesystem containing all bundled .flf font files.
//
//go:embed fonts/*.flf
var FontFS embed.FS
