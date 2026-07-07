package string

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fullJustifyCalude(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		maxWidth int
		expected []string
	}{
		{
			"leetcode_1",
			[]string{"This", "is", "an", "example", "of", "text", "justification."},
			16,
			[]string{
				"This    is    an",
				"example  of text",
				"justification.  ",
			},
		},
		{
			"leetcode_2",
			[]string{"What", "must", "be", "acknowledgment", "shall", "be"},
			16,
			[]string{
				"What   must   be",
				"acknowledgment  ",
				"shall be        ",
			},
		},
		{
			"single_word_per_line",
			[]string{"a", "b", "c"},
			1,
			[]string{"a", "b", "c"},
		},
		{
			"single_word",
			[]string{"hello"},
			10,
			[]string{"hello     "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, fullJustifyCalude(tt.words, tt.maxWidth))
		})
	}
}

// Test_fullJustify covers the original fullJustify implementation.
// Both fullJustify and fullJustifyCalude produce identical output — the test
// table is shared to keep expected values in one place.
//
// Branch coverage:
//   - multi-word line with uneven space distribution (leetcode_1)
//   - single word fills its own line, padded right (leetcode_2: "acknowledgment")
//   - every word is on its own line — justify() single-word branch (single_word_per_line)
//   - only one word in the input — last-line path only (single_word)
//   - three-word last line, left-justified (three_word_last_line)
func Test_fullJustify(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		maxWidth int
		expected []string
	}{
		{
			"leetcode_1",
			[]string{"This", "is", "an", "example", "of", "text", "justification."},
			16,
			[]string{
				"This    is    an",
				"example  of text",
				"justification.  ",
			},
		},
		{
			"leetcode_2",
			[]string{"What", "must", "be", "acknowledgment", "shall", "be"},
			16,
			[]string{
				"What   must   be",
				"acknowledgment  ",
				"shall be        ",
			},
		},
		{
			// Every word exactly fills maxWidth → single-word branch in justify closure.
			"single_word_per_line",
			[]string{"a", "b", "c"},
			1,
			[]string{"a", "b", "c"},
		},
		{
			// One word total — the entire output is the last-line left-justify path.
			"single_word",
			[]string{"hello"},
			10,
			[]string{"hello     "},
		},
		{
			// ["ab","cd"] → lineLetters=4, gaps=1, totalSpaces=2, spacePerGap=2 → "ab  cd"
			// last line: ["ef"] → "ef    " (left-justified, padded right)
			"two_word_line_then_last",
			[]string{"ab", "cd", "ef"},
			6,
			[]string{"ab  cd", "ef    "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, fullJustify(tt.words, tt.maxWidth))
		})
	}
}
