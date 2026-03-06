package figlet

import "strings"

// ---------------------------------------------------------------------------
// Horizontal smushing rules (hRule1–6)
// Each function returns (result string, ok bool).  ok == false means the
// two characters cannot be smushed by that rule.
// ---------------------------------------------------------------------------

// hRule1Smush — Equal Character Smushing (code value 1)
// Two identical sub-characters smush into one, except for hardblanks.
func hRule1Smush(ch1, ch2, hardBlank string) (string, bool) {
	if ch1 == ch2 && ch1 != hardBlank {
		return ch1, true
	}

	return "", false
}

// hRule2Smush — Underscore Smushing (code value 2)
// An underscore is replaced by any of: | / \ [ ] { } ( ) < >
func hRule2Smush(ch1, ch2 string) (string, bool) {
	const rule2Str = `|/\[]{}()<>`

	if ch1 == "_" && strings.ContainsRune(rule2Str, rune(ch2[0])) {
		return ch2, true
	}

	if ch2 == "_" && strings.ContainsRune(rule2Str, rune(ch1[0])) {
		return ch1, true
	}

	return "", false
}

// hRule3Smush — Hierarchy Smushing (code value 4)
// Class hierarchy: | > /\ > [] > {} > () > <>
// When two chars are from different classes, the one from the latter wins.
func hRule3Smush(ch1, ch2 string) (string, bool) {
	const rule3Classes = "| /\\ [] {} () <>"

	p1 := strings.Index(rule3Classes, ch1)
	p2 := strings.Index(rule3Classes, ch2)

	if p1 != -1 && p2 != -1 {
		diff := p1 - p2
		if diff < 0 {
			diff = -diff
		}

		if p1 != p2 && diff != 1 {
			start := p1
			if p2 > p1 {
				start = p2
			}

			return string(rule3Classes[start]), true
		}
	}

	return "", false
}

// hRule4Smush — Opposite Pair Smushing (code value 8)
// Opposing brackets/braces/parens smush into |
func hRule4Smush(ch1, ch2 string) (string, bool) {
	const rule4Str = "[] {} ()"

	p1 := strings.Index(rule4Str, ch1)
	p2 := strings.Index(rule4Str, ch2)

	if p1 != -1 && p2 != -1 {
		diff := p1 - p2
		if diff < 0 {
			diff = -diff
		}

		if diff <= 1 {
			return "|", true
		}
	}

	return "", false
}

// hRule5Smush — Big X Smushing (code value 16)
// /\ → |,  \/ → Y,  >< → X
func hRule5Smush(ch1, ch2 string) (string, bool) {
	patterns := map[string]string{
		`/\`: "|",
		`\/`: "Y",
		`><`: "X",
	}

	if v, ok := patterns[ch1+ch2]; ok {
		return v, true
	}

	return "", false
}

// hRule6Smush — Hardblank Smushing (code value 32)
// Two hardblanks smush into one hardblank.
func hRule6Smush(ch1, ch2, hardBlank string) (string, bool) {
	if ch1 == hardBlank && ch2 == hardBlank {
		return hardBlank, true
	}

	return "", false
}

// ---------------------------------------------------------------------------
// Vertical smushing rules (vRule1–5)
// ---------------------------------------------------------------------------

// vRule1Smush — Equal Character Smushing (code value 256)
func vRule1Smush(ch1, ch2 string) (string, bool) {
	if ch1 == ch2 {
		return ch1, true
	}

	return "", false
}

// vRule2Smush — Underscore Smushing (code value 512) — delegates to hRule2
func vRule2Smush(ch1, ch2 string) (string, bool) {
	return hRule2Smush(ch1, ch2)
}

// vRule3Smush — Hierarchy Smushing (code value 1024) — delegates to hRule3
func vRule3Smush(ch1, ch2 string) (string, bool) {
	return hRule3Smush(ch1, ch2)
}

// vRule4Smush — Horizontal Line Smushing (code value 2048)
// Stacked "-" and "_" (in either order) produce "=".
func vRule4Smush(ch1, ch2 string) (string, bool) {
	if (ch1 == "-" && ch2 == "_") || (ch1 == "_" && ch2 == "-") {
		return "=", true
	}

	return "", false
}

// vRule5Smush — Vertical Line Supersmushing (code value 4096)
// Two vertical bars smush into one.
func vRule5Smush(ch1, ch2 string) (string, bool) {
	if ch1 == "|" && ch2 == "|" {
		return "|", true
	}

	return "", false
}

// ---------------------------------------------------------------------------
// Universal smush
// ---------------------------------------------------------------------------

// uniSmush overrides ch1 with ch2, treating space as transparent.
// ch2 is space  → return ch1
// ch2 is hardblank and ch1 is not space → return ch1
// otherwise     → return ch2
func uniSmush(ch1, ch2, hardBlank string) string {
	if ch2 == " " || ch2 == "" {
		return ch1
	}

	if ch2 == hardBlank && ch1 != " " {
		return ch1
	}

	return ch2
}

// ---------------------------------------------------------------------------
// Vertical smush helpers
// ---------------------------------------------------------------------------

// canVerticalSmush reports whether two lines of FIGlet art can be
// vertically smushed given the current options.
// Returns "valid", "end", or "invalid".
func canVerticalSmush(txt1, txt2 string, opts InternalOptions) string {
	if opts.FittingRules.VLayout == lFullWidth {
		return "invalid"
	}

	minLen := len(txt1)
	if len(txt2) < minLen {
		minLen = len(txt2)
	}

	if minLen == 0 {
		return "invalid"
	}

	endSmush := false

	for ii := 0; ii < minLen; ii++ {
		ch1 := string(txt1[ii])
		ch2 := string(txt2[ii])

		if ch1 != " " && ch2 != " " {
			switch opts.FittingRules.VLayout {
			case lFitting:
				return "invalid"
			case lSmushing:
				return "end"
			default: // lControlledSmushing
				if _, ok := vRule5Smush(ch1, ch2); ok {
					// super-smushing: continue but don't yet mark endSmush
					continue
				}

				validSmush := false
				if !validSmush && opts.FittingRules.VRule1 {
					_, validSmush = vRule1Smush(ch1, ch2)
				}

				if !validSmush && opts.FittingRules.VRule2 {
					_, validSmush = vRule2Smush(ch1, ch2)
				}

				if !validSmush && opts.FittingRules.VRule3 {
					_, validSmush = vRule3Smush(ch1, ch2)
				}

				if !validSmush && opts.FittingRules.VRule4 {
					_, validSmush = vRule4Smush(ch1, ch2)
				}

				endSmush = true
				if !validSmush {
					return "invalid"
				}
			}
		}
	}

	if endSmush {
		return "end"
	}

	return "valid"
}

// getVerticalSmushDist returns the number of rows by which two blocks of
// FIGlet art can overlap vertically.
func getVerticalSmushDist(lines1, lines2 []string, opts InternalOptions) int {
	maxDist := len(lines1)
	len1 := len(lines1)
	curDist := 1
	result := ""

	for curDist <= maxDist {
		start1 := len1 - curDist
		if start1 < 0 {
			start1 = 0
		}

		subLines1 := lines1[start1:len1]
		end2 := curDist

		if end2 > maxDist {
			end2 = maxDist
		}

		if end2 > len(lines2) {
			end2 = len(lines2)
		}

		subLines2 := lines2[:end2]

		slen := len(subLines2)
		result = ""

		for ii := range slen {
			ret := canVerticalSmush(subLines1[ii], subLines2[ii], opts)

			if ret == "end" {
				result = ret
			} else if ret == "invalid" {
				result = ret
				break
			} else {
				if result == "" {
					result = "valid"
				}
			}
		}

		if result == "invalid" {
			curDist--
			break
		}

		if result == "end" {
			break
		}

		if result == "valid" {
			curDist++
		}
	}

	if maxDist < curDist {
		return maxDist
	}

	return curDist
}

// verticallySmushLines merges two lines of FIGlet art row by row.
func verticallySmushLines(line1, line2 string, opts InternalOptions) string {
	minLen := len(line1)
	if len(line2) < minLen {
		minLen = len(line2)
	}

	var sb strings.Builder
	fr := opts.FittingRules

	for ii := 0; ii < minLen; ii++ {
		ch1 := string(line1[ii])
		ch2 := string(line2[ii])

		if ch1 != " " && ch2 != " " {
			switch fr.VLayout {
			case lFitting, lSmushing:
				sb.WriteString(uniSmush(ch1, ch2, opts.HardBlank))
			default: // lControlledSmushing
				smushed := ""
				ok := false
				if !ok && fr.VRule5 {
					smushed, ok = vRule5Smush(ch1, ch2)
				}

				if !ok && fr.VRule1 {
					smushed, ok = vRule1Smush(ch1, ch2)
				}

				if !ok && fr.VRule2 {
					smushed, ok = vRule2Smush(ch1, ch2)
				}

				if !ok && fr.VRule3 {
					smushed, ok = vRule3Smush(ch1, ch2)
				}

				if !ok && fr.VRule4 {
					smushed, ok = vRule4Smush(ch1, ch2)
				}

				if ok {
					sb.WriteString(smushed)
				} else {
					sb.WriteString(uniSmush(ch1, ch2, opts.HardBlank))
				}
			}
		} else {
			sb.WriteString(uniSmush(ch1, ch2, opts.HardBlank))
		}
	}

	return sb.String()
}

// verticalSmush merges two FIGlet text blocks with the given vertical overlap.
func verticalSmush(lines1, lines2 []string, overlap int, opts InternalOptions) []string {
	len1 := len(lines1)
	len2 := len(lines2)

	// piece1: rows of lines1 before the overlap zone
	piece1Start := 0
	piece1End := len1 - overlap
	if piece1End < 0 {
		piece1End = 0
	}

	piece1 := lines1[piece1Start:piece1End]

	// piece2: merged overlap rows
	over1Start := len1 - overlap
	if over1Start < 0 {
		over1Start = 0
	}

	piece2_1 := lines1[over1Start:len1]
	over2End := overlap
	if over2End > len2 {
		over2End = len2
	}

	piece2_2 := lines2[:over2End]

	piece2 := make([]string, len(piece2_1))
	for ii := range piece2_1 {
		if ii >= len2 {
			piece2[ii] = piece2_1[ii]
		} else {
			piece2[ii] = verticallySmushLines(piece2_1[ii], piece2_2[ii], opts)
		}
	}

	// piece3: remaining rows of lines2 after the overlap zone
	piece3Start := overlap
	if piece3Start > len2 {
		piece3Start = len2
	}

	piece3 := lines2[piece3Start:]

	result := make([]string, 0, len(piece1)+len(piece2)+len(piece3))
	result = append(result, piece1...)
	result = append(result, piece2...)
	result = append(result, piece3...)
	return result
}

// getHorizontalSmushLength returns how many columns two FIGlet art rows can
// horizontally overlap.
func getHorizontalSmushLength(txt1, txt2 string, opts InternalOptions) int {
	fr := opts.FittingRules
	if fr.HLayout == lFullWidth {
		return 0
	}

	len1 := len(txt1)
	len2 := len(txt2)
	maxDist := len1
	curDist := 1
	breakAfter := false

	if len1 == 0 {
		return 0
	}

distCal:
	for curDist <= maxDist {
		seg1Start := len1 - curDist
		seg1 := txt1[seg1Start : seg1Start+curDist]
		end2 := curDist
		if end2 > len2 {
			end2 = len2
		}

		seg2 := txt2[:end2]

		for ii := 0; ii < end2; ii++ {
			ch1 := string(seg1[ii])
			ch2 := string(seg2[ii])

			if ch1 != " " && ch2 != " " {
				switch fr.HLayout {
				case lFitting:
					curDist--
					break distCal
				case lSmushing:
					if ch1 == opts.HardBlank || ch2 == opts.HardBlank {
						curDist-- // universal smushing does not smush hardblanks
					}

					break distCal
				default: // lControlledSmushing
					breakAfter = true
					validSmush := false
					if !validSmush && fr.HRule1 {
						_, validSmush = hRule1Smush(ch1, ch2, opts.HardBlank)
					}

					if !validSmush && fr.HRule2 {
						_, validSmush = hRule2Smush(ch1, ch2)
					}

					if !validSmush && fr.HRule3 {
						_, validSmush = hRule3Smush(ch1, ch2)
					}

					if !validSmush && fr.HRule4 {
						_, validSmush = hRule4Smush(ch1, ch2)
					}

					if !validSmush && fr.HRule5 {
						_, validSmush = hRule5Smush(ch1, ch2)
					}

					if !validSmush && fr.HRule6 {
						_, validSmush = hRule6Smush(ch1, ch2, opts.HardBlank)
					}

					if !validSmush {
						curDist--
						break distCal
					}
				}
			}
		}

		if breakAfter {
			break
		}

		curDist++
	}

	if maxDist < curDist {
		return maxDist
	}

	return curDist
}

// horizontalSmush merges two FIGlet character blocks with the given overlap.
func horizontalSmush(block1, block2 []string, overlap int, opts InternalOptions) []string {
	fr := opts.FittingRules
	height := opts.Height
	out := make([]string, height)

	for ii := range height {
		txt1 := block1[ii]
		txt2 := block2[ii]
		len1 := len(txt1)
		len2 := len(txt2)

		overlapStart := len1 - overlap
		if overlapStart < 0 {
			overlapStart = 0
		}

		piece1 := txt1[:overlapStart]

		// Overlap segment from block1 and block2
		seg1Start := len1 - overlap
		if seg1Start < 0 {
			seg1Start = 0
		}

		seg1 := txt1[seg1Start:]
		end2 := overlap
		if end2 > len2 {
			end2 = len2
		}

		seg2 := txt2[:end2]

		var piece2 strings.Builder
		for jj := range overlap {
			ch1 := " "
			if jj < len1 && jj < len(seg1) {
				ch1 = string(seg1[jj])
			}

			ch2 := " "
			if jj < len2 && jj < len(seg2) {
				ch2 = string(seg2[jj])
			}

			if ch1 != " " && ch2 != " " {
				if fr.HLayout == lFitting || fr.HLayout == lSmushing {
					piece2.WriteString(uniSmush(ch1, ch2, opts.HardBlank))
				} else {
					// Controlled smushing
					nextCh := ""
					ok := false
					if !ok && fr.HRule1 {
						nextCh, ok = hRule1Smush(ch1, ch2, opts.HardBlank)
					}

					if !ok && fr.HRule2 {
						nextCh, ok = hRule2Smush(ch1, ch2)
					}

					if !ok && fr.HRule3 {
						nextCh, ok = hRule3Smush(ch1, ch2)
					}

					if !ok && fr.HRule4 {
						nextCh, ok = hRule4Smush(ch1, ch2)
					}

					if !ok && fr.HRule5 {
						nextCh, ok = hRule5Smush(ch1, ch2)
					}

					if !ok && fr.HRule6 {
						nextCh, ok = hRule6Smush(ch1, ch2, opts.HardBlank)
					}

					if !ok {
						nextCh = uniSmush(ch1, ch2, opts.HardBlank)
					}

					piece2.WriteString(nextCh)
				}
			} else {
				piece2.WriteString(uniSmush(ch1, ch2, opts.HardBlank))
			}
		}

		piece3 := ""
		if overlap < len2 {
			piece3 = txt2[overlap:]
		}

		out[ii] = piece1 + piece2.String() + piece3
	}

	return out
}
