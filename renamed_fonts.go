package figlet

var renamedFonts = map[string]string{
	"ANSI-Compact":     "ANSI Compact",
	"Patorjk's Cheese": "Patorjks Cheese",
}

func getFontName(name string) string {
	if v, ok := renamedFonts[name]; ok {
		return v
	}

	return name
}
