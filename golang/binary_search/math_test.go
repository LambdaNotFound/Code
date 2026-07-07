package binarysearch

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_myPow(t *testing.T) {
	testCases := []struct {
		name     string
		x        float64
		n        int
		expected float64
	}{
		{"positive_exponent", 2.0, 10, 1024.0},
		{"zero_exponent", 2.0, 0, 1.0},
		{"one_exponent", 5.0, 1, 5.0},
		{"negative_exponent", 2.0, -2, 0.25},
		{"base_one", 1.0, 1000, 1.0},
		{"base_zero", 0.0, 5, 0.0},
		{"fractional_base", 0.5, 3, 0.125},
		{"negative_base_even", -2.0, 4, 16.0},
		{"negative_base_odd", -2.0, 3, -8.0},
		{"large_exponent", 1.00001, 100000, math.Pow(1.00001, 100000)},
	}

	const epsilon = 1e-9
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := myPow(tc.x, tc.n)
			assert.InDelta(t, tc.expected, got, epsilon)

			got = myPowRecursive(tc.x, tc.n)
			assert.InDelta(t, tc.expected, got, epsilon)
		})
	}
}
