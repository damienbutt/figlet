package figlet

// PrintDirection controls the direction in which FIGlet characters are rendered.
type PrintDirection int

const (
	DefaultDirection PrintDirection = iota - 1 // -1: use font header's value
	LeftToRight                                //  0: left to right (FIGlet default)
	RightToLeft                                //  1: right to left
)

// FittingRules holds the horizontal and vertical layout mode and the
// individual smushing rule flags decoded from a font's header.
type FittingRules struct {
	HLayout int
	HRule1  bool
	HRule2  bool
	HRule3  bool
	HRule4  bool
	HRule5  bool
	HRule6  bool
	VLayout int
	VRule1  bool
	VRule2  bool
	VRule3  bool
	VRule4  bool
	VRule5  bool
}

// FontMetadata contains the parsed header fields of a FIGlet font file (.flf).
type FontMetadata struct {
	Baseline        int
	CodeTagCount    *int
	FittingRules    FittingRules
	FullLayout      *int
	HardBlank       string
	Height          int
	MaxLength       int
	NumCommentLines int
	OldLayout       int
	PrintDirection  PrintDirection
}

// FontName is a string alias for a FIGlet font name.
// KerningMethods is a string alias for a horizontal kerning/smushing mode.
// FittingProperties is a string alias for a field name within FittingRules.
type (
	FontName          string
	KerningMethods    string
	FittingProperties string
)

// KerningMethods constants define the supported horizontal kerning/smushing modes.
const (
	KerningDefault            KerningMethods = "default"
	KerningFull               KerningMethods = "full"
	KerningFitted             KerningMethods = "fitted"
	KerningControlledSmushing KerningMethods = "controlled smushing"
	KerningUniversalSmushing  KerningMethods = "universal smushing"
)

// FittingProperties constants name the individual fields of FittingRules.
const (
	FitHLayout FittingProperties = "hLayout"
	FitHRule1  FittingProperties = "hRule1"
	FitHRule2  FittingProperties = "hRule2"
	FitHRule3  FittingProperties = "hRule3"
	FitHRule4  FittingProperties = "hRule4"
	FitHRule5  FittingProperties = "hRule5"
	FitHRule6  FittingProperties = "hRule6"
	FitVLayout FittingProperties = "vLayout"
	FitVRule1  FittingProperties = "vRule1"
	FitVRule2  FittingProperties = "vRule2"
	FitVRule3  FittingProperties = "vRule3"
	FitVRule4  FittingProperties = "vRule4"
	FitVRule5  FittingProperties = "vRule5"
)

// FigletOptions controls how a piece of text is rendered by Text.
// All fields are optional; zero values fall back to font/package defaults.
type FigletOptions struct {
	Font             FontName
	HorizontalLayout KerningMethods
	VerticalLayout   KerningMethods
	Width            int
	WhitespaceBreak  bool
	PrintDirection   PrintDirection
	ShowHardBlanks   bool
}

// FigletDefaults holds the package-level defaults used when no FigletOptions
// are provided to Text. Modify via Defaults().
type FigletDefaults struct {
	Font               FontName
	FontPath           string
	FetchFontIfMissing bool
}

// FigletFont is the internal representation of a parsed .flf font file.
type FigletFont struct {
	options  *FontMetadata
	comment  string
	numChars int
	charCode map[int][]string
}

// NewFigletFont returns an empty FigletFont ready to be populated by ParseFont.
func NewFigletFont() *FigletFont {
	return &FigletFont{
		options:  &FontMetadata{},
		comment:  "",
		numChars: 0,
		charCode: make(map[int][]string),
	}
}

type InternalOptions struct {
	FontMetadata
	Width           int
	WhitespaceBreak bool
	ShowHardBlanks  bool
}

type FigCharWithOverlap struct {
	fig     []string
	overlap int
}

type FigCharsWithOverlap struct {
	chars   []FigCharWithOverlap
	overlap int
}

type BreakWordResult struct {
	outputFigText []string
	chars         []FigCharWithOverlap
}
