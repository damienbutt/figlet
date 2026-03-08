package figlet_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/damienbutt/figlet"
)

// readExpected reads a golden file from testdata/expected/ and strips any trailing newline.
func readExpected(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile("testdata/expected/" + filename)
	require.NoError(t, err)
	return strings.TrimSuffix(string(data), "\n")
}

func getMaxWidth(input string) int {
	max := 0
	for line := range strings.SplitSeq(input, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}

	return max
}

// standardMeta reflects the known metadata of the Standard.flf font.
var standardMeta = figlet.FontMetadata{
	HardBlank:       "$",
	Height:          6,
	Baseline:        5,
	MaxLength:       16,
	OldLayout:       15,
	NumCommentLines: 13,
	PrintDirection:  figlet.LeftToRight,
	FullLayout:      new(24463),
	CodeTagCount:    new(229),
	FittingRules: figlet.FittingRules{
		VLayout: 3,
		VRule5:  true,
		VRule4:  true,
		VRule3:  true,
		VRule2:  true,
		VRule1:  true,
		HLayout: 3,
		HRule6:  false,
		HRule5:  false,
		HRule4:  true,
		HRule3:  true,
		HRule2:  true,
		HRule1:  true,
	},
}

// ---------------------------------------------------------------------------
// Standard font
// ---------------------------------------------------------------------------

func TestStandardFont(t *testing.T) {
	t.Run("renders text with standard font and default vertical layout", func(t *testing.T) {
		actual, err := figlet.Text("FIGlet\nFonts", &figlet.FigletOptions{
			Font: "Standard",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "standard_default"), actual)
	})

	t.Run("renders text with standard font and fitted vertical layout", func(t *testing.T) {
		actual, err := figlet.Text("FIGlet\nFONTS", &figlet.FigletOptions{
			Font:           "Standard",
			VerticalLayout: "fitted",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "standard"), actual)
	})

	t.Run("renders text with a parsed font", func(t *testing.T) {
		data, err := os.ReadFile("fonts/Standard.flf")
		require.NoError(t, err)

		_, err = figlet.ParseFont("StandardParseFontName", string(data))
		require.NoError(t, err)

		actual, err := figlet.Text("FIGlet\nFONTS", &figlet.FigletOptions{
			Font:           "StandardParseFontName",
			VerticalLayout: "fitted",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "standard"), actual)
	})
}

// ---------------------------------------------------------------------------
// Graffiti font
// ---------------------------------------------------------------------------

func TestGraffitiFont(t *testing.T) {
	t.Run("renders text with graffiti font and fitted horizontal layout", func(t *testing.T) {
		actual, err := figlet.Text("ABC.123", &figlet.FigletOptions{
			Font:             "Graffiti",
			HorizontalLayout: "fitted",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "graffiti"), actual)
	})
}

// ---------------------------------------------------------------------------
// Text wrapping
// ---------------------------------------------------------------------------

func TestTextWrapping(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		opts     figlet.FigletOptions
		expected string
		maxWidth int
	}{
		{
			name:     "wrap simple",
			text:     "Hello From The Figlet Library",
			opts:     figlet.FigletOptions{Font: "Standard", Width: 80},
			expected: "wrapSimple",
			maxWidth: 80,
		},
		{
			name:     "wrap word",
			text:     "Hello From The Figlet Library",
			opts:     figlet.FigletOptions{Font: "Standard", Width: 80, WhitespaceBreak: true},
			expected: "wrapWord",
			maxWidth: 80,
		},
		{
			name:     "wrap simple three lines",
			text:     "Hello From The Figlet Library That Wrap Text",
			opts:     figlet.FigletOptions{Font: "Standard", Width: 80},
			expected: "wrapSimpleThreeLines",
			maxWidth: 80,
		},
		{
			name:     "wrap word three lines",
			text:     "Hello From The Figlet Library That Wrap Text",
			opts:     figlet.FigletOptions{Font: "Standard", Width: 80, WhitespaceBreak: true},
			expected: "wrapWordThreeLines",
			maxWidth: 80,
		},
		{
			name:     "wrap whitespace break word",
			text:     "Hello LongLongLong Word Longerhello",
			opts:     figlet.FigletOptions{Font: "Standard", Width: 30, WhitespaceBreak: true},
			expected: "wrapWhitespaceBreakWord",
			maxWidth: 30,
		},
		{
			name:     "wrap whitespace log string",
			text:     "xxxxxxxxxxxxxxxxxxxxxxxx",
			opts:     figlet.FigletOptions{Font: "Standard", Width: 30, WhitespaceBreak: true},
			expected: "wrapWhitespaceLogString",
			maxWidth: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := figlet.Text(tt.text, &tt.opts)
			require.NoError(t, err)
			assert.LessOrEqual(t, getMaxWidth(actual), tt.maxWidth)
			assert.Equal(t, readExpected(t, tt.expected), actual)
		})
	}
}

// ---------------------------------------------------------------------------
// Misc fonts
// ---------------------------------------------------------------------------

func TestMiscFonts(t *testing.T) {
	t.Run("renders with Dancing Font using full horizontal layout", func(t *testing.T) {
		actual, err := figlet.Text("pizzapie", &figlet.FigletOptions{
			Font:             "Dancing Font",
			HorizontalLayout: "full",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "dancingFont"), actual)
	})

	t.Run("renders with Dancing Font right-to-left", func(t *testing.T) {
		actual, err := figlet.Text("pizzapie", &figlet.FigletOptions{
			Font:             "Dancing Font",
			HorizontalLayout: "full",
			PrintDirection:   figlet.RightToLeft,
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "dancingFontReverse"), actual)
	})

	t.Run("follows vertical smush rule 2 (Slant font)", func(t *testing.T) {
		actual, err := figlet.Text("Terminal\nChess", &figlet.FigletOptions{
			Font: "Slant",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "verticalSmushRule2"), actual)
	})

	t.Run("multiline text with empty lines (miniwi font)", func(t *testing.T) {
		actual, err := figlet.Text("This\n\nis\n\n\na test", &figlet.FigletOptions{
			Font: "miniwi",
		})

		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "miniwi_multiline"), actual)
	})
}

// ---------------------------------------------------------------------------
// Fonts list
// ---------------------------------------------------------------------------

func TestFonts(t *testing.T) {
	t.Run("returns a non-empty list containing known fonts", func(t *testing.T) {
		fonts, err := figlet.Fonts()
		require.NoError(t, err)
		assert.NotEmpty(t, fonts)
		assert.Contains(t, fonts, "Standard")
		assert.Contains(t, fonts, "Graffiti")
	})

	t.Run("all fonts load and produce non-empty output", func(t *testing.T) {
		fonts, err := figlet.Fonts()
		require.NoError(t, err)

		for _, font := range fonts {
			font := font
			t.Run(font, func(t *testing.T) {
				actual, err := figlet.Text("abc ABC ...", &figlet.FigletOptions{
					Font: figlet.FontName(font),
				})

				require.NoError(t, err)
				assert.Greater(t, getMaxWidth(actual), 0)
			})
		}
	})
}

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

func TestErrorHandling(t *testing.T) {
	t.Run("returns error containing 'Font' when font not found", func(t *testing.T) {
		_, err := figlet.Text("test", &figlet.FigletOptions{Font: "NonExistentFont"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Font")
	})
}

// ---------------------------------------------------------------------------
// Renamed fonts
// ---------------------------------------------------------------------------

func TestRenamedFonts(t *testing.T) {
	t.Run("ANSI-Compact (hyphen variant) resolves correctly", func(t *testing.T) {
		actual, err := figlet.Text("this is a test", &figlet.FigletOptions{Font: "ANSI-Compact"})
		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "ansiCompact"), actual)

		_, comment, err := figlet.Metadata("ANSI-Compact")
		require.NoError(t, err)
		assert.Equal(t, 13, strings.Index(comment, "Loic"))
	})

	t.Run("ANSI Compact (space variant) resolves correctly", func(t *testing.T) {
		actual, err := figlet.Text("this is a test", &figlet.FigletOptions{Font: "ANSI Compact"})
		require.NoError(t, err)
		assert.Equal(t, readExpected(t, "ansiCompact"), actual)

		_, comment, err := figlet.Metadata("ANSI Compact")
		require.NoError(t, err)
		assert.Equal(t, 13, strings.Index(comment, "Loic"))
	})
}

// ---------------------------------------------------------------------------
// Font loading and metadata
// ---------------------------------------------------------------------------

func TestLoadFont(t *testing.T) {
	t.Run("loads Standard-Test font from testdata and returns correct metadata", func(t *testing.T) {
		figlet.ClearLoadedFonts()
		defer figlet.ClearLoadedFonts()

		// Point to local testdata font directory
		original := figlet.Defaults(nil)
		figlet.Defaults(&figlet.FigletDefaults{FontPath: "testdata/fonts"})
		defer figlet.Defaults(&figlet.FigletDefaults{FontPath: original.FontPath})

		assert.Empty(t, figlet.LoadedFonts())

		meta, err := figlet.LoadFont("Standard-Test")
		require.NoError(t, err)
		assert.Equal(t, standardMeta, *meta)
		assert.Contains(t, figlet.LoadedFonts(), "Standard-Test")
	})
}
