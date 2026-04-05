package solid_coding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_myAtoi(t *testing.T) {
	testCases := []struct {
		name     string
		s        string
		expected int
	}{
		{
			"case 1",
			"42",
			42,
		},
		{
			"case 2",
			"   -042",
			-42,
		},
		{
			"case 3",
			"1337c0d3",
			1337,
		},
		{
			"case 4",
			"0-1",
			0,
		},
		{
			"case 5",
			"words and 987",
			0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := myAtoi(tc.s)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_longestPalindromeLength(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			"case 1",
			"abccccdd",
			7,
		},
		{
			"case 2",
			"a",
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := longestPalindromeLength(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
