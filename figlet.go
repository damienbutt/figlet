package figlet

import "errors"

var errNotImplemented = errors.New("not implemented")

// TextSync generates ASCII art from the given text using the provided options.
func TextSync(text string, opts *FigletOptions) (string, error) {
	return "", errNotImplemented
}

// LoadFont loads and returns the metadata for a named font.
func LoadFont(name string) (*FontMetadata, error) {
	return nil, errNotImplemented
}

// Fonts returns a sorted list of all available font names.
func Fonts() ([]string, error) {
	return nil, errNotImplemented
}
