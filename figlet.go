package figlet

import "errors"

var errNotImplemented = errors.New("not implemented")

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
	return "", errNotImplemented
}

// Text is an alias for TextSync.
func Text(text string, opts *FigletOptions) (string, error) {
	return TextSync(text, opts)
}

// LoadFont loads a font by name from the configured font path or embedded fonts.
func LoadFont(name string) (*FontMetadata, error) {
	return nil, errNotImplemented
}

// ParseFont parses raw .flf font data and registers it under the given name.
func ParseFont(name string, data string) (*FontMetadata, error) {
	return nil, errNotImplemented
}

// PreloadFonts loads multiple fonts into the cache.
func PreloadFonts(names []FontName) error {
	return errNotImplemented
}

// LoadedFonts returns the names of all currently cached fonts.
func LoadedFonts() []string {
	return nil
}

// ClearLoadedFonts removes all fonts from the cache.
func ClearLoadedFonts() {}

// Fonts returns a sorted list of all available font names.
func Fonts() ([]string, error) {
	return nil, errNotImplemented
}

// Metadata returns the metadata and comment string for a named font.
func Metadata(name string) (*FontMetadata, string, error) {
	return nil, "", errNotImplemented
}
