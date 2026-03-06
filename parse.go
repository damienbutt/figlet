package figlet

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// escapeRegexChar escapes a single character that may have special meaning in
// a regular expression.
func escapeRegexChar(char string) string {
	if strings.ContainsAny(char, `.*+?^${}()|[\]`) {
		return `\` + char
	}

	return char
}

// removeEndChar strips the FIGlet end marker(s) from a font line.
// The last non-space character of a trimmed line is the end character.
// The bottom line (lineNum == fontHeight-1) may have one or two end chars;
// all other lines have exactly one.
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

// ParseFont parses raw .flf font data and registers it into the cache under
// the given name. If a font with that name is already cached it is
// overwritten.
func ParseFont(name string, data string) (*FontMetadata, error) {
	// Normalise line endings; strip BOM if present
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.ReplaceAll(data, "\r", "\n")
	data = strings.TrimPrefix(data, "\xef\xbb\xbf") // UTF-8 BOM

	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("Font file is empty")
	}

	// --- Parse header ---
	headerLine := lines[0]
	lines = lines[1:]

	headerData := strings.Fields(headerLine)
	if len(headerData) < 6 {
		return nil, fmt.Errorf("Font file header is too short")
	}

	if len(headerData[0]) < 6 {
		return nil, fmt.Errorf("Font file has invalid header")
	}

	hardBlank := string(headerData[0][5])

	parseInt := func(s string) (int, error) {
		return strconv.Atoi(s)
	}

	height, err := parseInt(headerData[1])
	if err != nil {
		return nil, fmt.Errorf("Font file: invalid height: %w", err)
	}

	baseline, err := parseInt(headerData[2])
	if err != nil {
		return nil, fmt.Errorf("Font file: invalid baseline: %w", err)
	}

	maxLength, err := parseInt(headerData[3])
	if err != nil {
		return nil, fmt.Errorf("Font file: invalid maxLength: %w", err)
	}

	oldLayout, err := parseInt(headerData[4])
	if err != nil {
		return nil, fmt.Errorf("Font file: invalid oldLayout: %w", err)
	}

	numCommentLines, err := parseInt(headerData[5])
	if err != nil {
		return nil, fmt.Errorf("Font file: invalid numCommentLines: %w", err)
	}

	// Validate required fields
	if utf8.RuneCountInString(hardBlank) != 1 {
		return nil, fmt.Errorf("Font file: invalid hardBlank character")
	}

	printDirection := PrintDirection(0) // default: left-to-right
	if len(headerData) >= 7 {
		pd, err := parseInt(headerData[6])
		if err != nil {
			return nil, fmt.Errorf("Font file: invalid printDirection: %w", err)
		}

		printDirection = PrintDirection(pd)
	} else {
		printDirection = LeftToRight
	}

	var fullLayout *int
	if len(headerData) >= 8 {
		fl, err := parseInt(headerData[7])
		if err != nil {
			return nil, fmt.Errorf("Font file: invalid fullLayout: %w", err)
		}

		fullLayout = &fl
	}

	var codeTagCount *int
	if len(headerData) >= 9 {
		ct, err := parseInt(headerData[8])
		if err != nil {
			return nil, fmt.Errorf("Font file: invalid codeTagCount: %w", err)
		}

		codeTagCount = &ct
	}

	opts := FontMetadata{
		HardBlank:       hardBlank,
		Height:          height,
		Baseline:        baseline,
		MaxLength:       maxLength,
		OldLayout:       oldLayout,
		NumCommentLines: numCommentLines,
		PrintDirection:  printDirection,
		FullLayout:      fullLayout,
		CodeTagCount:    codeTagCount,
	}

	opts.FittingRules = getSmushingRules(oldLayout, fullLayout)

	// --- Required char codes: 32–126, then the extended Latin chars ---
	charNums := make([]int, 0, 104)
	for i := 32; i <= 126; i++ {
		charNums = append(charNums, i)
	}

	charNums = append(charNums, 196, 214, 220, 228, 246, 252, 223)

	// Validate we have enough data
	needed := numCommentLines + height*len(charNums)
	if len(lines) < needed {
		return nil, fmt.Errorf(
			"Font file is missing data. Line length: %d. Comment lines: %d. Height: %d. Num chars: %d",
			len(lines), numCommentLines, height, len(charNums),
		)
	}

	font := NewFigletFont()
	font.options = &opts

	// --- Parse comment ---
	font.comment = strings.Join(lines[:numCommentLines], "\n")
	lines = lines[numCommentLines:]

	// --- Parse required characters ---
	for _, cNum := range charNums {
		charLines := lines[:height]
		lines = lines[height:]

		for i := range height {
			if i >= len(charLines) {
				charLines = append(charLines, "")
			}

			charLines[i] = removeEndChar(charLines[i], i, height)
		}

		font.charCode[cNum] = charLines
		font.numChars++
	}

	// --- Parse optional tagged characters ---
	for len(lines) > 0 {
		cNumLine := lines[0]
		lines = lines[1:]

		cNumLine = strings.TrimSpace(cNumLine)
		if cNumLine == "" {
			break
		}

		parts := strings.SplitN(cNumLine, " ", 2)
		raw := parts[0]

		var parsedNum int
		switch {
		case regexp.MustCompile(`^-?0[xX][0-9a-fA-F]+$`).MatchString(raw):
			n, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimPrefix(raw, "-"), "0x"), 16, 64)
			if err != nil {
				n, err = strconv.ParseInt(strings.TrimPrefix(strings.TrimPrefix(raw, "-"), "0X"), 16, 64)
			}

			if err != nil {
				return nil, fmt.Errorf("Font file: error parsing data. Invalid data: %s", raw)
			}

			if strings.HasPrefix(raw, "-") {
				n = -n
			}

			parsedNum = int(n)
		case regexp.MustCompile(`^-?0[0-7]+$`).MatchString(raw):
			n, err := strconv.ParseInt(strings.TrimPrefix(raw, "-"), 8, 64)
			if err != nil {
				return nil, fmt.Errorf("Font file: error parsing data. Invalid data: %s", raw)
			}

			if strings.HasPrefix(raw, "-") {
				n = -n
			}

			parsedNum = int(n)
		case regexp.MustCompile(`^-?[0-9]+$`).MatchString(raw):
			parsedNum, err = parseInt(raw)
			if err != nil {
				return nil, fmt.Errorf("Font file: error parsing data. Invalid data: %s", raw)
			}
		default:
			return nil, fmt.Errorf("Font file: error parsing data. Invalid data: %s", raw)
		}

		// Per FIGlet spec: code must be in range [-2147483648, 2147483647], excluding -1.
		if parsedNum == -1 {
			return nil, fmt.Errorf("Font file: error parsing data. The char code -1 is not permitted.")
		}

		if parsedNum < -2147483648 || parsedNum > 2147483647 {
			if parsedNum < -2147483648 {
				return nil, fmt.Errorf("Font file: error parsing data. The char code cannot be less than -2147483648.")
			}

			return nil, fmt.Errorf("Font file: error parsing data. The char code cannot be greater than 2147483647.")
		}

		if len(lines) < height {
			break
		}

		charLines := lines[:height]
		lines = lines[height:]

		for i := range height {
			charLines[i] = removeEndChar(charLines[i], i, height)
		}

		font.charCode[parsedNum] = charLines
		font.numChars++
	}

	// Store in package-level cache
	storeFigFont(name, font)

	return font.options, nil
}
