package main

import (
	"sort"

	figlet "github.com/damienbutt/figlet"
)

// figletFonts returns a sorted list of all available embedded font names.
func figletFonts() ([]string, error) {
	entries, err := figlet.FontFS.ReadDir("fonts")
	if err != nil {
		return nil, err
	}

	var fonts []string
	for _, e := range entries {
		name := e.Name()
		if len(name) > 4 && name[len(name)-4:] == ".flf" {
			fonts = append(fonts, name[:len(name)-4])
		}
	}

	sort.Strings(fonts)
	return fonts, nil
}

// figletLoadFont loads and returns the metadata for a named font.
func figletLoadFont(name string) (*figlet.FontMetadata, error) {
	return figlet.LoadFont(name)
}

// figletTextSync generates ASCII art from text using the given options.
func figletTextSync(text, font, horizontalLayout, verticalLayout string, width int) (string, error) {
	opts := &figlet.FigletOptions{
		Font:             figlet.FontName(font),
		HorizontalLayout: figlet.KerningMethods(horizontalLayout),
		VerticalLayout:   figlet.KerningMethods(verticalLayout),
		Width:            width,
	}

	return figlet.TextSync(text, opts)
}
