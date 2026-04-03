package dynamic_programming

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_numDecodings(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected int
	}{
		{name: "example1", s: "12", expected: 2},           // "AB" or "L"
		{name: "example2", s: "226", expected: 3},          // "BZ", "VF", "BBF"
		{name: "leading_zero", s: "06", expected: 0},       // 0 can't map to any letter
		{name: "single_digit", s: "1", expected: 1},
		{name: "single_nine", s: "9", expected: 1},
		{name: "double_zero", s: "00", expected: 0},
		{name: "ten", s: "10", expected: 1},                // only "J"
		{name: "twenty_six", s: "26", expected: 2},         // "BF" or "Z"
		{name: "twenty_seven", s: "27", expected: 1},       // only "BG" (27>26)
		{name: "long_string", s: "11106", expected: 2},     // "AAJF" or "KJF"
		{name: "all_ones", s: "1111", expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, numDecodings(tt.s))
			assert.Equal(t, tt.expected, numDecodingsTopDown(tt.s))
			assert.Equal(t, tt.expected, numDecodingsTopDownMemo(tt.s))
		})
	}
}
