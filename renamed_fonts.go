package figlet

// renamedFonts maps the canonical display name (as users / callers know it)
// to the actual filename stem on disk / in the embedded FS.
// Reasons for renaming:
//   - ANSI-Compact: hyphen variant used by callers; file is "ANSI Compact.flf"
//   - Patorjk's Cheese: apostrophe is invalid in go:embed patterns; file is "Patorjks Cheese.flf"
var renamedFonts = map[string]string{
	"ANSI-Compact":     "ANSI Compact",
	"Patorjk's Cheese": "Patorjks Cheese",
}

// diskToDisplay is the inverse of renamedFonts: disk stem → display name.
// Used by Fonts() so the listing shows the canonical names callers expect.
var diskToDisplay map[string]string

func init() {
	diskToDisplay = make(map[string]string, len(renamedFonts))
	for display, disk := range renamedFonts {
		diskToDisplay[disk] = display
	}
}

// getFontName resolves a caller-supplied name to the actual filename stem.
func getFontName(name string) string {
	if v, ok := renamedFonts[name]; ok {
		return v
	}

	return name
}

// getDisplayName returns the display name for a disk stem, falling back to
// the stem itself when no alias exists.
func getDisplayName(diskStem string) string {
	if v, ok := diskToDisplay[diskStem]; ok {
		return v
	}

	return diskStem
}
