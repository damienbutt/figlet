package figlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHRule1Smush(t *testing.T) {
	tests := []struct {
		name      string
		ch1, ch2  string
		hardBlank string
		want      string
		ok        bool
	}{
		{"identical chars", "a", "a", "$", "a", true},
		{"identical spaces", " ", " ", "$", " ", true},
		{"identical — both are hardblank → false", "$", "$", "$", "", false},
		{"different chars", "a", "b", "$", "", false},
		{"ch1 is hardblank, ch2 is not", "$", "a", "$", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hRule1Smush(tt.ch1, tt.ch2, tt.hardBlank)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHRule2Smush(t *testing.T) {
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"underscore left, pipe right", "_", "|", "|", true},
		{"underscore left, slash right", "_", "/", "/", true},
		{"underscore left, backslash right", "_", `\`, `\`, true},
		{"underscore left, left-bracket right", "_", "[", "[", true},
		{"pipe left, underscore right", "|", "_", "|", true},
		{"slash left, underscore right", "/", "_", "/", true},
		{"unrelated chars", "a", "b", "", false},
		{"both underscores", "_", "_", "", false},
		{"underscore left, non-member right", "_", "a", "", false},
		{"non-member left, underscore right", "a", "_", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hRule2Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHRule3Smush(t *testing.T) {
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"pipe vs bracket — bracket wins", "|", "[", "[", true},
		{"bracket vs pipe — bracket wins", "[", "|", "[", true},
		{"bracket vs brace — brace wins", "[", "{", "{", true},
		{"pipe vs paren — paren wins", "|", "(", "(", true},
		{"pipe vs angle — angle wins", "|", "<", "<", true},
		// Adjacent chars belong to the same class (diff==1) → no smush
		{"adjacent: slash and backslash (same /\\ class)", "/", `\`, "", false},
		{"adjacent: open/close bracket", "[", "]", "", false},
		{"adjacent: open/close brace", "{", "}", "", false},
		// Chars not in the hierarchy at all
		{"char not in hierarchy", "a", "|", "", false},
		{"both chars not in hierarchy", "a", "b", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hRule3Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHRule4Smush(t *testing.T) {
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"open/close bracket → |", "[", "]", "|", true},
		{"close/open bracket → |", "]", "[", "|", true},
		{"open/close brace → |", "{", "}", "|", true},
		{"close/open brace → |", "}", "{", "|", true},
		{"open/close paren → |", "(", ")", "|", true},
		{"bracket vs brace (different pairs)", "[", "}", "", false},
		{"bracket vs paren (different pairs)", "[", ")", "", false},
		{"unrelated chars", "a", "b", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hRule4Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHRule5Smush(t *testing.T) {
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"/\\ → |", "/", `\`, "|", true},
		{"\\/ → Y", `\`, "/", "Y", true},
		{">< → X", ">", "<", "X", true},
		{"<> → no match (reversed)", "<", ">", "", false},
		{"unrelated pair", "a", "b", "", false},
		{"single valid char only", "/", "a", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hRule5Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHRule6Smush(t *testing.T) {
	tests := []struct {
		name      string
		ch1, ch2  string
		hardBlank string
		want      string
		ok        bool
	}{
		{"both are hardblank → hardblank", "$", "$", "$", "$", true},
		{"custom hardblank char", "^", "^", "^", "^", true},
		{"ch1 hardblank, ch2 not", "$", "a", "$", "", false},
		{"ch2 hardblank, ch1 not", "a", "$", "$", "", false},
		{"neither is hardblank", "a", "a", "$", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hRule6Smush(tt.ch1, tt.ch2, tt.hardBlank)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVRule1Smush(t *testing.T) {
	// vRule1 is equal-char smushing with no hardblank exception (unlike hRule1)
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"identical chars", "a", "a", "a", true},
		{"identical spaces", " ", " ", " ", true},
		{"identical — even dollar sign (no hardblank exception)", "$", "$", "$", true},
		{"different chars", "a", "b", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := vRule1Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVRule2Smush(t *testing.T) {
	// vRule2 delegates to hRule2 — verify delegation with one case each direction
	t.Run("underscore left, pipe right", func(t *testing.T) {
		got, ok := vRule2Smush("_", "|")
		assert.True(t, ok)
		assert.Equal(t, "|", got)
	})

	t.Run("unrelated chars → false", func(t *testing.T) {
		_, ok := vRule2Smush("a", "b")
		assert.False(t, ok)
	})
}

func TestVRule3Smush(t *testing.T) {
	// vRule3 delegates to hRule3 — verify delegation with one case each direction
	t.Run("pipe vs bracket — bracket wins", func(t *testing.T) {
		got, ok := vRule3Smush("|", "[")
		assert.True(t, ok)
		assert.Equal(t, "[", got)
	})

	t.Run("adjacent chars → false", func(t *testing.T) {
		_, ok := vRule3Smush("[", "]")
		assert.False(t, ok)
	})
}

func TestVRule4Smush(t *testing.T) {
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"dash over underscore → =", "-", "_", "=", true},
		{"underscore over dash → =", "_", "-", "=", true},
		{"both dashes", "-", "-", "", false},
		{"both underscores", "_", "_", "", false},
		{"unrelated", "a", "b", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := vRule4Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVRule5Smush(t *testing.T) {
	tests := []struct {
		name     string
		ch1, ch2 string
		want     string
		ok       bool
	}{
		{"both pipes → |", "|", "|", "|", true},
		{"pipe and space", "|", " ", "", false},
		{"space and pipe", " ", "|", "", false},
		{"unrelated chars", "a", "b", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := vRule5Smush(tt.ch1, tt.ch2)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUniSmush(t *testing.T) {
	tests := []struct {
		name      string
		ch1, ch2  string
		hardBlank string
		want      string
	}{
		{"ch2 is space → ch1 wins", "a", " ", "$", "a"},
		{"ch2 is empty → ch1 wins", "a", "", "$", "a"},
		{"ch2 is hardblank and ch1 is not space → ch1 wins", "a", "$", "$", "a"},
		{"ch2 is hardblank and ch1 is space → ch2 wins (hardblank)", " ", "$", "$", "$"},
		{"ch2 is regular char → ch2 wins", "a", "b", "$", "b"},
		{"both space → space (ch2 wins)", " ", " ", "$", " "},
		{"ch1 space, ch2 regular → ch2 wins", " ", "x", "$", "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniSmush(tt.ch1, tt.ch2, tt.hardBlank)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCanVerticalSmush(t *testing.T) {
	makeOpts := func(vLayout int, r1, r2, r3, r4, r5 bool) InternalOptions {
		return InternalOptions{
			FontMetadata: FontMetadata{
				HardBlank: "$",
				FittingRules: FittingRules{
					VLayout: vLayout,
					VRule1:  r1,
					VRule2:  r2,
					VRule3:  r3,
					VRule4:  r4,
					VRule5:  r5,
				},
			},
		}
	}

	tests := []struct {
		name       string
		txt1, txt2 string
		opts       InternalOptions
		want       string
	}{
		{
			"fullWidth layout → always invalid",
			"abc", "def",
			makeOpts(lFullWidth, false, false, false, false, false),
			"invalid",
		},
		{
			"empty strings → invalid (minLen is 0)",
			"", "",
			makeOpts(lFitting, false, false, false, false, false),
			"invalid",
		},
		{
			"all spaces → valid (loop never enters non-space branch)",
			"   ", "   ",
			makeOpts(lFitting, false, false, false, false, false),
			"valid",
		},
		{
			"non-space, fitting layout → invalid",
			"a", "b",
			makeOpts(lFitting, false, false, false, false, false),
			"invalid",
		},
		{
			"non-space, universal smushing → end",
			"a", "b",
			makeOpts(lSmushing, false, false, false, false, false),
			"end",
		},
		{
			"controlled smushing, vRule5 (both pipes) → valid (super-smush, loop continues)",
			"|", "|",
			makeOpts(lControlledSmushing, false, false, false, false, true),
			"valid",
		},
		{
			"controlled smushing, vRule1 matches (same char) → end",
			"a", "a",
			makeOpts(lControlledSmushing, true, false, false, false, false),
			"end",
		},
		{
			"controlled smushing, no rule matches → invalid",
			"a", "b",
			makeOpts(lControlledSmushing, false, false, false, false, false),
			"invalid",
		},
		{
			"one side is space → valid (space is transparent)",
			"a", " ",
			makeOpts(lFitting, false, false, false, false, false),
			"valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canVerticalSmush(tt.txt1, tt.txt2, tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}
