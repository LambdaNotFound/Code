package solid_coding

import "strings"

/*
 * 68. Text Justification
 *
 * How do you greedily pack words onto each line?
 * Given totalSpaces and gaps, how do you compute spacesPerGap and extraSpaces?
 * How do you handle the single word edge case on a line — no gaps to distribute spaces into?
 * How do you differentiate the last line from regular lines?
 *
 */
func fullJustifyCalude(words []string, maxWidth int) []string {
	res := make([]string, 0)
	line := make([]string, 0)
	lineLetters := 0

	for _, word := range words {
		// +len(line) accounts for minimum 1 space between each word
		if lineLetters+len(line)+len(word) > maxWidth {
			res = append(res, justify(line, lineLetters, maxWidth))
			line = line[:0]
			lineLetters = 0
		}
		line = append(line, word)
		lineLetters += len(word)
	}

	// last line: left justified — single space between words, pad right
	lastLine := strings.Join(line, " ")
	lastLine += strings.Repeat(" ", maxWidth-len(lastLine))
	res = append(res, lastLine)

	return res
}

func justify(line []string, lineLetters, maxWidth int) string {
	if len(line) == 1 {
		// single word — pad right with spaces
		return line[0] + strings.Repeat(" ", maxWidth-lineLetters)
	}

	totalSpaces := maxWidth - lineLetters
	gaps := len(line) - 1
	spacePerGap := totalSpaces / gaps
	extraSpaces := totalSpaces % gaps

	var sb strings.Builder
	for i, word := range line {
		sb.WriteString(word)
		if i < gaps { // distribute extra spaces to leftmost gaps first
			spaces := spacePerGap
			if i < extraSpaces {
				spaces++
			}
			sb.WriteString(strings.Repeat(" ", spaces))
		}
	}

	return sb.String()
}
