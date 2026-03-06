package main

import (
	figlet "github.com/damienbutt/figlet"
)

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
