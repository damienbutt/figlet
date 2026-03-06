package figlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// newFigChar
// ---------------------------------------------------------------------------

func TestNewFigChar(t *testing.T) {
	t.Run("height 0 returns empty (non-nil) slice", func(t *testing.T) {
		assert.Equal(t, []string{}, newFigChar(0))
	})

	t.Run("returns slice of empty strings with correct length", func(t *testing.T) {
		assert.Equal(t, []string{"", "", ""}, newFigChar(3))
	})

	t.Run("height 1 returns single-element slice", func(t *testing.T) {
		assert.Equal(t, []string{""}, newFigChar(1))
	})
}

// ---------------------------------------------------------------------------
// figLinesWidth
// ---------------------------------------------------------------------------

func TestFigLinesWidth(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{"empty slice", []string{}, 0},
		{"single line", []string{"hello"}, 5},
		{"longest line wins", []string{"hello", "world!"}, 6},
		{"all empty strings", []string{"", ""}, 0},
		{"unicode rune count not byte count", []string{"日本語"}, 3},
		{"mixed ascii and unicode", []string{"hi", "日本語"}, 3},
		{"empty string mixed with spaces", []string{"", "  "}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, figLinesWidth(tt.lines))
		})
	}
}

// ---------------------------------------------------------------------------
// padLines
// ---------------------------------------------------------------------------

func TestPadLines(t *testing.T) {
	t.Run("pads each line by n spaces", func(t *testing.T) {
		assert.Equal(t, []string{"a  ", "bc  "}, padLines([]string{"a", "bc"}, 2))
	})

	t.Run("zero padding returns lines unchanged", func(t *testing.T) {
		assert.Equal(t, []string{"a"}, padLines([]string{"a"}, 0))
	})

	t.Run("empty slice returns empty (non-nil) slice", func(t *testing.T) {
		assert.Equal(t, []string{}, padLines([]string{}, 3))
	})

	t.Run("returns a new slice — original is not modified", func(t *testing.T) {
		original := []string{"a", "b"}
		out := padLines(original, 1)
		assert.Equal(t, []string{"a ", "b "}, out)
		assert.Equal(t, []string{"a", "b"}, original)
	})
}

// ---------------------------------------------------------------------------
// getFontName (display/alias → disk stem)
// ---------------------------------------------------------------------------

func TestGetFontName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"ANSI-Compact alias resolves to disk name", "ANSI-Compact", "ANSI Compact"},
		{"Patorjk's Cheese resolves to sanitised disk name", "Patorjk's Cheese", "Patorjks Cheese"},
		{"unknown name returns unchanged", "Standard", "Standard"},
		{"already-resolved name returns unchanged", "ANSI Compact", "ANSI Compact"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getFontName(tt.input))
		})
	}
}

// ---------------------------------------------------------------------------
// getDisplayName (disk stem → display/alias)
// ---------------------------------------------------------------------------

func TestGetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		diskStem string
		want     string
	}{
		{"Patorjks Cheese maps back to apostrophe name", "Patorjks Cheese", "Patorjk's Cheese"},
		{"ANSI Compact maps back to hyphen alias", "ANSI Compact", "ANSI-Compact"},
		{"unknown disk stem returned unchanged", "Standard", "Standard"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getDisplayName(tt.diskStem))
		})
	}
}

// ---------------------------------------------------------------------------
// Round-trip: every renamedFonts entry must survive display→disk→display
// ---------------------------------------------------------------------------

func TestRenamedFontsRoundTrip(t *testing.T) {
	for display, disk := range renamedFonts {
		t.Run(display, func(t *testing.T) {
			assert.Equal(t, disk, getFontName(display), "getFontName should resolve display→disk")
			assert.Equal(t, display, getDisplayName(disk), "getDisplayName should resolve disk→display")
		})
	}
}
