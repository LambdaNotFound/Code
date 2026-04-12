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
