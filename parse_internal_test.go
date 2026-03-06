package figlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// escapeRegexChar
// ---------------------------------------------------------------------------

func TestEscapeRegexChar(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// All 14 regex-special characters must be escaped.
		{".", `\.`},
		{"*", `\*`},
		{"+", `\+`},
		{"?", `\?`},
		{"^", `\^`},
		{"$", `\$`},
		{"{", `\{`},
		{"}", `\}`},
		{"(", `\(`},
		{")", `\)`},
		{"|", `\|`},
		{"[", `\[`},
		{`\`, `\\`},
		{"]", `\]`},
		// Non-special characters must pass through unchanged.
		{"@", "@"},
		{"a", "a"},
		{" ", " "},
		{"0", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, escapeRegexChar(tt.input))
		})
	}
}

// ---------------------------------------------------------------------------
// removeEndChar
// ---------------------------------------------------------------------------

func TestRemoveEndChar(t *testing.T) {
	tests := []struct {
		name       string
		line       string
		lineNum    int
		fontHeight int
		want       string
	}{
		{
			"non-bottom line strips one end char",
			"  abc@  ", 0, 4,
			"  abc",
		},
		{
			"bottom line strips two consecutive end chars",
			"  abc@@  ", 3, 4,
			"  abc",
		},
		{
			"bottom line with a single end char is also valid",
			"  abc@  ", 3, 4,
			"  abc",
		},
		{
			"end char derived from last non-space char (pipe)",
			"  abc|", 0, 4,
			"  abc",
		},
		{
			"bottom line, pipe end char",
			"  abc||", 3, 4,
			"  abc",
		},
		{
			"line containing only the end marker → empty string",
			"@", 0, 4,
			"",
		},
		{
			"empty input → empty output",
			"", 0, 4,
			"",
		},
		{
			"regex-special end char (+) on non-bottom line",
			"abc+", 0, 4,
			"abc",
		},
		{
			"regex-special end char (+) on bottom line",
			"abc++", 3, 4,
			"abc",
		},
		{
			"height-1 font: line 0 is the bottom line",
			"abc@@", 0, 1,
			"abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeEndChar(tt.line, tt.lineNum, tt.fontHeight)
			assert.Equal(t, tt.want, got)
		})
	}
}
