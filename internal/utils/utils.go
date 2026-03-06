package utils

import (
	"regexp"
	"strings"
)

// Escapes chars that have special meaning in Regex
func escapeRegexChar(char string) string {
	if strings.ContainsAny(char, `.*+?^${}()|[\]`) {
		return "\\" + char
	}

	return char
}

// Figures out the end char for a FIGlet line and removes it. Technically there aren't supposed to be white spaces
// after the end char, but certain TOIlet fonts have this. The FIGlet unit app handles this though so we handle it
// here too.
func removeEndChar(line string, lineNum int, fontHeight int) string {
	endChar := "@"

	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 0 {
		endChar = escapeRegexChar(string(trimmed[len(trimmed)-1]))
	}

	var pattern string
	if lineNum == fontHeight-1 {
		pattern = endChar + endChar + `?\s*$`
	} else {
		pattern = endChar + `\s*$`
	}

	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(line, "")
}
