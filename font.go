package figlet

type PrintDirection int

const (
	DefaultDirection PrintDirection = iota - 1 // -1: use font header's value
	LeftToRight                                //  0: left to right (FIGlet default)
	RightToLeft                                //  1: right to left
)

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

type (
	FontName       string
	KerningMethods string
)

type FigletOptions struct {
	Font             FontName
	HorizontalLayout KerningMethods
	VerticalLayout   KerningMethods
	Width            int
	WhitespaceBreak  bool
	PrintDirection   PrintDirection
	ShowHardBlanks   bool
}

type FigletDefaults struct {
	Font               FontName
	FontPath           string
	FetchFontIfMissing bool
}

type FigletFont struct {
	options  *FontMetadata
	comment  string
	numChars int
	charCode map[int][]string
}

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
