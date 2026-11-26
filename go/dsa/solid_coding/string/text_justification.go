package string

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
func fullJustify(words []string, maxWidth int) []string {
	res := []string{}
	lineWords := []string{}
	lineLetters := 0

	justify := func(lineWords []string, lineLetters int) string {
		if len(lineWords) == 1 { // handle single word in a line
			return lineWords[0] + strings.Repeat(" ", maxWidth-lineLetters)
		}

		gaps := len(lineWords) - 1
		spacePerGap := (maxWidth - lineLetters) / gaps
		spaceExtra := (maxWidth - lineLetters) % gaps

		var sb strings.Builder
		for i, word := range lineWords {
			sb.WriteString(word)
			if i < gaps {
				spaces := spacePerGap
				if i < spaceExtra {
					spaces++
				}
				sb.WriteString(strings.Repeat(" ", spaces))
			}
		}
		return sb.String()
	}

	for _, word := range words {
		if lineLetters+len(word)+len(lineWords) > maxWidth {
			line := justify(lineWords, lineLetters)
			res = append(res, line)

			lineWords = []string{}
			lineLetters = 0
		}

		lineWords = append(lineWords, word)
		lineLetters = lineLetters + len(word)
	}

	lastLine := strings.Join(lineWords, " ") // justify last line
	lastLine = lastLine + strings.Repeat(" ", maxWidth-len(lastLine))
	res = append(res, lastLine)

	return res
}

func fullJustifyCalude(words []string, maxWidth int) []string {
	res := make([]string, 0)
	line := make([]string, 0)
	lineLetters := 0

	for _, word := range words {
		// +len(line) accounts for minimum 1 space between each word
		if lineLetters+len(line)+len(word) > maxWidth {
			res = append(res, justify(line, lineLetters, maxWidth))
			line = line[:0] // reset slice
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
