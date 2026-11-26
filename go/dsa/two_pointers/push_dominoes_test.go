package two_pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pushDominoes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"RR.L", "RR.L"},        // R...L meeting in middle with even gap
		{".L.R...LR..L..", "LL.RR.LLRRLL.."}, // LC example 2
		{"R.R.L", "RRR.L"},     // chained R push
		{"L.R", "L.R"},          // L...R stays neutral
		{"R.L", "R.L"},           // span=1, R...L: 0R + 1dot + 0L, middle stays
		{"R..L", "RRLL"},        // span=2, R...L: 1R + 1L, no middle dot
		{"R...L", "RR.LL"},      // span=3, R...L: 1R + 1dot + 1L
		{".....", "....."},       // all standing
		{"LLLLL", "LLLLL"},      // all L
		{"RRRRR", "RRRRR"},      // all R
		{"R", "R"},               // single R
		{"L", "L"},               // single L
		{".", "."},               // single standing
		{"LR", "LR"},             // adjacent L R, nothing moves
		{"RL", "RL"},             // adjacent R L, already fallen
		{"R....L", "RRRLLL"},    // span=4, R...L: 2R + 2L, anchor chars included
	}

	fns := []struct {
		name string
		fn   func(string) string
	}{
		{"pushDominoes", pushDominoes},
		{"pushDominoesVisualized", pushDominoesVisualized},
		{"pushDominoesBFS", pushDominoesBFS},
	}

	for _, f := range fns {
		for _, tt := range tests {
			t.Run(f.name+"/"+tt.input, func(t *testing.T) {
				assert.Equal(t, tt.expected, f.fn(tt.input))
			})
		}
	}
}
