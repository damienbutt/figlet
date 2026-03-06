package figlet_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/damienbutt/figlet"
)

// minimalFontData returns a syntactically valid .flf data string with the
// minimum content needed to pass ParseFont (height=1, no comments, blank glyphs).
func minimalFontData() string {
	var sb strings.Builder
	sb.WriteString("flf2a$ 1 1 1 0 0\n")
	// 95 required chars (32–126) + 7 extended Latin = 102 total
	for range 102 {
		sb.WriteString("@\n")
	}

	return sb.String()
}

// fontWithTaggedChar appends a tagged-character section to a minimal valid font.
func fontWithTaggedChar(tagLine string) string {
	return minimalFontData() + tagLine + "\n@\n"
}

// ---------------------------------------------------------------------------
// ParseFont error paths
// ---------------------------------------------------------------------------

func TestParseFontErrors(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr string
	}{
		{
			"empty input",
			"",
			"header is too short",
		},
		{
			"header fields too few",
			"flf2a$\n",
			"header is too short",
		},
		{
			"header first field too short (< 6 bytes)",
			"fl 1 1 1 0 0\n",
			"has invalid header",
		},
		{
			"invalid height",
			"flf2a$ X 1 1 0 0\n",
			"invalid height",
		},
		{
			"invalid baseline",
			"flf2a$ 1 X 1 0 0\n",
			"invalid baseline",
		},
		{
			"invalid maxLength",
			"flf2a$ 1 1 X 0 0\n",
			"invalid maxLength",
		},
		{
			"invalid oldLayout",
			"flf2a$ 1 1 1 X 0\n",
			"invalid oldLayout",
		},
		{
			"invalid numCommentLines",
			"flf2a$ 1 1 1 0 X\n",
			"invalid numCommentLines",
		},
		{
			"header only — no char data",
			"flf2a$ 1 1 1 0 0\n",
			"missing data",
		},
		{
			"tagged char with code -1",
			fontWithTaggedChar("-1"),
			"char code -1 is not permitted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := figlet.ParseFont("__test__"+tt.name, tt.data)
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

// ---------------------------------------------------------------------------
// PreloadFonts
// ---------------------------------------------------------------------------

func TestPreloadFonts(t *testing.T) {
	t.Run("successfully preloads multiple fonts into cache", func(t *testing.T) {
		figlet.ClearLoadedFonts()
		defer figlet.ClearLoadedFonts()

		err := figlet.PreloadFonts([]figlet.FontName{"Standard", "Graffiti"})
		require.NoError(t, err)

		loaded := figlet.LoadedFonts()
		assert.Contains(t, loaded, "Standard")
		assert.Contains(t, loaded, "Graffiti")
	})

	t.Run("returns error when font does not exist", func(t *testing.T) {
		err := figlet.PreloadFonts([]figlet.FontName{"NonExistentFont"})
		require.Error(t, err)
	})

	t.Run("empty list returns no error", func(t *testing.T) {
		err := figlet.PreloadFonts([]figlet.FontName{})
		require.NoError(t, err)
	})
}

// ---------------------------------------------------------------------------
// Fonts() display names
// ---------------------------------------------------------------------------

func TestFontsDisplayNames(t *testing.T) {
	fonts, err := figlet.Fonts()
	require.NoError(t, err)

	assert.Contains(t, fonts, "Patorjk's Cheese", "apostrophe must be restored from disk stem")
	assert.NotContains(t, fonts, "Patorjks Cheese", "sanitised disk name must not appear in listing")

	// ANSI-Compact: renamedFonts maps the hyphen-alias → disk name, so the
	// listing shows the hyphen-alias (the full inversion of renamedFonts).
	assert.Contains(t, fonts, "ANSI-Compact", "hyphen alias must appear in font listing")
	assert.NotContains(t, fonts, "ANSI Compact", "raw disk name must not appear in listing")
}
