package figlet

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

var figDefaults = FigletDefaults{
	Font:     "Standard",
	FontPath: "fonts",
}

// Defaults gets or sets the default options. Pass nil to read without modifying.
func Defaults(opts *FigletDefaults) FigletDefaults {
	if opts != nil {
		if opts.Font != "" {
			figDefaults.Font = opts.Font
		}

		if opts.FontPath != "" {
			figDefaults.FontPath = opts.FontPath
		}
	}

	return figDefaults
}

// TextSync generates ASCII art from the given text using the provided options.
func TextSync(text string, opts *FigletOptions) (string, error) {
	fontName := figDefaults.Font
	if opts != nil && opts.Font != "" {
		fontName = FontName(opts.Font)
	}

	meta, err := LoadFont(string(fontName))
	if err != nil {
		return "", err
	}

	internalOpts := reworkFontOpts(*meta, opts)
	return generateText(string(fontName), internalOpts, text)
}

// Text is an alias for TextSync.
func Text(text string, opts *FigletOptions) (string, error) {
	return TextSync(text, opts)
}

// PreloadFonts loads multiple fonts into the cache.
func PreloadFonts(names []FontName) error {
	for _, name := range names {
		if _, err := LoadFont(string(name)); err != nil {
			return err
		}
	}

	return nil
}

// Fonts returns a sorted list of all available font names (from embedded FS).
func Fonts() ([]string, error) {
	entries, err := fs.ReadDir(FontFS, "fonts")
	if err != nil {
		return nil, fmt.Errorf("Font list unavailable: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".flf") {
			names = append(names, strings.TrimSuffix(e.Name(), ".flf"))
		}
	}

	sort.Strings(names)
	return names, nil
}

// Metadata returns the metadata and comment string for a named font.
func Metadata(name string) (*FontMetadata, string, error) {
	meta, err := LoadFont(name)
	if err != nil {
		return nil, "", err
	}

	actualName := getFontName(name)
	mu.RLock()
	font, ok := figFonts[actualName]
	mu.RUnlock()

	comment := ""
	if ok {
		comment = font.comment
	}

	return meta, comment, nil
}
