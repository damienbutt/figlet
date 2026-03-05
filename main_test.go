package main

import "testing"

func TestEscapeRegexChar(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
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
		{"]", `\]`},
	}

	for _, test := range tests {
		result := escapeRegexChar(test.input)
		if result != test.expected {
			t.Errorf("escapeRegexChar(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestRemoveEndChar(t *testing.T) {
	tests := []struct {
		line       string
		lineNum    int
		fontHeight int
		expected   string
	}{
		{"Hello@", 0, 5, "Hello"},
		{"World@@", 4, 5, "World"},
		{"Test@ ", 2, 5, "Test"},
		{"Example@", 1, 5, "Example"},
	}

	for _, test := range tests {
		result := removeEndChar(test.line, test.lineNum, test.fontHeight)
		if result != test.expected {
			t.Errorf("removeEndChar(%q, %d, %d) = %q; expected %q", test.line, test.lineNum, test.fontHeight, result, test.expected)
		}
	}
}
