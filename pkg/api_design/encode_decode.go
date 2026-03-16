package apidesign

import (
	"strconv"
	"strings"
)

/**
 * Encode and Decode Strings
 *
 */
type Solution struct{}

func (s *Solution) Encode(strs []string) string {
	list := make([]string, 0)
	for _, str := range strs {
		l := len(str)
		s := strconv.Itoa(l) + "#" + str
		list = append(list, s)
	}

	return strings.Join(list, "")
}

func (s *Solution) EncodeClaude(strs []string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(strconv.Itoa(len(str)))
		sb.WriteByte('#')
		sb.WriteString(str)
	}
	return sb.String()
}

func (s *Solution) Decode(encoded string) []string {
	list := make([]string, 0)

	for left, right := 0, 0; right < len(encoded); {
		if encoded[right] != '#' {
			right++
		} else {
			l, _ := strconv.Atoi(encoded[left:right])
			str := encoded[right+1 : right+1+l]
			list = append(list, str)

			left = right + 1 + l
			right = left
		}
	}

	return list
}

func (s *Solution) DecodeClaude(encoded string) []string {
	result := make([]string, 0)
	left := 0

	for left < len(encoded) {
		// find delimiter '#'
		right := left
		for encoded[right] != '#' {
			right++
		}

		// extract length prefix and jump to string content
		l, _ := strconv.Atoi(encoded[left:right])
		str := encoded[right+1 : right+1+l]
		result = append(result, str)
		left = right + 1 + l
	}

	return result
}
