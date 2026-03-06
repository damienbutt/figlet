package figlet

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// newFigChar returns a blank FIGlet character (height lines of empty strings).
func newFigChar(height int) []string {
	lines := make([]string, height)
	for i := range lines {
		lines[i] = ""
	}

	return lines
}

// figLinesWidth returns the visual width (rune count) of the longest line.
func figLinesWidth(lines []string) int {
	max := 0
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > max {
			max = n
		}
	}

	return max
}

// padLines appends n spaces to every line in the slice.
func padLines(lines []string, n int) []string {
	pad := strings.Repeat(" ", n)
	out := make([]string, len(lines))

	for i, l := range lines {
		out[i] = l + pad
	}

	return out
}

// joinFigArray reduces a list of FigCharWithOverlap values into a single
// FIGlet text block by repeatedly calling horizontalSmush.
func joinFigArray(array []FigCharWithOverlap, height int, opts InternalOptions) []string {
	acc := newFigChar(height)
	for _, data := range array {
		acc = horizontalSmush(acc, data.fig, data.overlap, opts)
	}

	return acc
}

// breakWord finds the longest prefix of figChars that fits within opts.Width
// and returns both the assembled line and the remaining characters.
func breakWord(figChars []FigCharWithOverlap, height int, opts InternalOptions) BreakWordResult {
	for i := len(figChars) - 1; i > 0; i-- {
		w := joinFigArray(figChars[:i], height, opts)
		if figLinesWidth(w) <= opts.Width {
			return BreakWordResult{
				outputFigText: w,
				chars:         figChars[i:],
			}
		}
	}

	return BreakWordResult{
		outputFigText: newFigChar(height),
		chars:         figChars,
	}
}

// smushVerticalFigLines merges an accumulated output block with an incoming
// block of FIGlet lines using vertical smushing.
func smushVerticalFigLines(output, lines []string, opts InternalOptions) []string {
	if len(output) == 0 || len(lines) == 0 {
		if len(output) == 0 {
			return lines
		}

		return output
	}

	len1 := figLinesWidth(output)
	len2 := figLinesWidth(lines)

	if len1 > len2 {
		lines = padLines(lines, len1-len2)
	} else if len2 > len1 {
		output = padLines(output, len2-len1)
	}

	overlap := getVerticalSmushDist(output, lines, opts)
	return verticalSmush(output, lines, overlap, opts)
}

// generateFigTextLines converts a single line of text into one or more rows
// of FIGlet art, applying wrapping when opts.Width > 0.
func generateFigTextLines(txt string, font *FigletFont, opts InternalOptions) [][]string {
	height := opts.Height
	outputFigLines := [][]string{}
	outputFigText := newFigChar(height)
	fr := opts.FittingRules

	// Right-to-left: reverse the string at the character level
	if opts.PrintDirection == RightToLeft {
		runes := []rune(txt)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}

		txt = string(runes)
	}

	type figCharsWithOverlap struct {
		chars   []FigCharWithOverlap
		overlap int
	}

	nextFigChars := figCharsWithOverlap{}
	var figWords []FigCharWithOverlap

	txtRunes := []rune(txt)
	length := len(txtRunes)

	for charIndex := range length {
		char := txtRunes[charIndex]
		isSpace := char == ' ' || char == '\t'
		figChar, hasFigChar := font.charCode[int(char)]

		if !hasFigChar {
			continue
		}

		overlap := 0
		if fr.HLayout != lFullWidth {
			overlap = 10000 // intentionally large; will be minimised

			for row := range height {
				ol := getHorizontalSmushLength(outputFigText[row], figChar[row], opts)
				if ol < overlap {
					overlap = ol
				}
			}

			if overlap == 10000 {
				overlap = 0
			}
		}

		if opts.Width > 0 {
			var textFigLine []string
			var maxWidth int

			if opts.WhitespaceBreak {
				// Next character added to the current word
				textFigWord := joinFigArray(
					append(nextFigChars.chars, FigCharWithOverlap{fig: figChar, overlap: overlap}),
					height, opts,
				)

				textFigLine = joinFigArray(
					append(figWords, FigCharWithOverlap{fig: textFigWord, overlap: nextFigChars.overlap}),
					height, opts,
				)

				maxWidth = figLinesWidth(textFigLine)
			} else {
				textFigLine = horizontalSmush(outputFigText, figChar, overlap, opts)
				maxWidth = figLinesWidth(textFigLine)
			}

			if maxWidth >= opts.Width && charIndex > 0 {
				if opts.WhitespaceBreak {
					// Emit everything up to (but not including) the last word
					outputFigText = joinFigArray(figWords[:max(0, len(figWords)-1)], height, opts)
					if len(figWords) > 1 {
						outputFigLines = append(outputFigLines, outputFigText)
						outputFigText = newFigChar(height)
					}

					figWords = nil
				} else {
					outputFigLines = append(outputFigLines, outputFigText)
					outputFigText = newFigChar(height)
				}
			}
		}

		if opts.Width > 0 && opts.WhitespaceBreak {
			if !isSpace || charIndex == length-1 {
				nextFigChars.chars = append(nextFigChars.chars, FigCharWithOverlap{fig: figChar, overlap: overlap})
			}

			if isSpace || charIndex == length-1 {
				// Break long words that still exceed the width
				var tmpBreak *BreakWordResult

				for {
					textFigLine := joinFigArray(nextFigChars.chars, height, opts)
					maxWidth := figLinesWidth(textFigLine)

					if maxWidth >= opts.Width {
						br := breakWord(nextFigChars.chars, height, opts)
						tmpBreak = &br
						nextFigChars = figCharsWithOverlap{chars: br.chars}
						outputFigLines = append(outputFigLines, br.outputFigText)
					} else {
						break
					}
				}

				// Add word to the line
				textFigLine := joinFigArray(nextFigChars.chars, height, opts)
				maxWidth := figLinesWidth(textFigLine)
				if maxWidth > 0 {
					if tmpBreak != nil {
						figWords = append(figWords, FigCharWithOverlap{fig: textFigLine, overlap: 1})
					} else {
						figWords = append(figWords, FigCharWithOverlap{fig: textFigLine, overlap: nextFigChars.overlap})
					}
				}

				// Save the space character itself for smushing
				if isSpace {
					figWords = append(figWords, FigCharWithOverlap{fig: figChar, overlap: overlap})
					outputFigText = newFigChar(height)
				}

				if charIndex == length-1 {
					// Last character: finalise the line
					outputFigText = joinFigArray(figWords, height, opts)
				}

				nextFigChars = figCharsWithOverlap{overlap: overlap}
				continue
			}
		}

		outputFigText = horizontalSmush(outputFigText, figChar, overlap, opts)
	}

	// Emit the final (possibly only) line if non-empty
	if figLinesWidth(outputFigText) > 0 {
		outputFigLines = append(outputFigLines, outputFigText)
	}

	// Replace hardblanks with spaces (unless showHardBlanks is set)
	if !opts.ShowHardBlanks && opts.HardBlank != "" {
		hb := opts.HardBlank
		for _, block := range outputFigLines {
			for row := range block {
				block[row] = strings.ReplaceAll(block[row], hb, " ")
			}
		}
	}

	// Special case: empty input produces one blank line of the right height
	if txt == "" && len(outputFigLines) == 0 {
		outputFigLines = append(outputFigLines, newFigChar(height))
	}

	return outputFigLines
}

// generateText renders txt into ASCII art using fontName with the provided
// InternalOptions, then returns the assembled string.
func generateText(fontName string, opts InternalOptions, txt string) (string, error) {
	txt = strings.ReplaceAll(txt, "\r\n", "\n")
	txt = strings.ReplaceAll(txt, "\r", "\n")

	actualFontName := getFontName(fontName)

	mu.RLock()
	font, ok := figFonts[actualFontName]
	mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("Font not loaded: %q", fontName)
	}

	inputLines := strings.Split(txt, "\n")
	var figLines [][]string

	for _, line := range inputLines {
		figLines = append(figLines, generateFigTextLines(line, font, opts)...)
	}

	if len(figLines) == 0 {
		return "", nil
	}

	output := figLines[0]
	for ii := 1; ii < len(figLines); ii++ {
		output = smushVerticalFigLines(output, figLines[ii], opts)
	}

	if output == nil {
		return "", nil
	}

	return strings.Join(output, "\n"), nil
}

// max returns the larger of a and b (for Go versions without built-in max).
// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }
