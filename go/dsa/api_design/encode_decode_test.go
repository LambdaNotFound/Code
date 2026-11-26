package apidesign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeDecodeStrings(t *testing.T) {
	testCases := []struct {
		name  string
		input []string
	}{
		{"basic", []string{"neet", "code", "love", "you"}},
		{"empty_strings", []string{"", ""}},
		{"single_word", []string{"hello"}},
		{"string_with_delimiter", []string{"a#b", "c#d"}}, // '#' inside string
		{"string_with_numbers_in_prefix", []string{"10hello", "world"}},
		{"empty_list", []string{}},
		{"single_empty_string", []string{""}},
		{"spaces_and_special", []string{"hello world", "foo#bar", "3#abc"}},
	}

	sol := &Solution{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := sol.Encode(tc.input)
			decoded := sol.Decode(encoded)
			assert.Equal(t, tc.input, decoded)

			// also verify the Claude variants round-trip
			encoded2 := sol.EncodeClaude(tc.input)
			decoded2 := sol.DecodeClaude(encoded2)
			assert.Equal(t, tc.input, decoded2)
		})
	}
}
