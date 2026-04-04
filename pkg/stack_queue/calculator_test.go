package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Basic arithmetic
		{"addition", "1 + 1", 2},
		{"subtraction", "2 - 1", 1},
		{"chain_add_sub", "2 - 1 + 2", 3},
		{"single_number", "42", 42},
		{"spaces_only_number", "   3   ", 3},

		// Parentheses
		{"single_parens", "(1 + 1)", 2},
		{"nested_parens", "(1+(4+5+2)-3)+(6+8)", 23},
		{"deep_nesting", "((((1+2))))", 3},
		{"double_parens", "((2))", 2},

		// Unary minus
		{"unary_before_parens", "-(3 + 4)", -7},
		{"unary_before_number", "-3 + 1", -2},
		{"unary_nested", "-(1-(2+3))", 4}, // -(1-5) = -(-4) = 4

		// Edge cases
		{"empty_string", "", 0},
		{"zero", "0", 0},
		{"negative_result", "1 - 5", -4},
		{"negative_inside_parens", "1 - (5 - 2)", -2},
		{"chain_parens", "(1)-(2)-(3)", -4},
		{"multiple_spaces", " 2-   1 +   2 ", 3},
		{"complex_spaces", " ( 7 -( 3 + (2 - 1) ) ) ", 3},
		{"large_number", "100 - 99", 1},
		{"multi_digit", "12 + 34", 46},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, calculate(tt.input))
		})
	}
}

func Test_calculateII(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// LeetCode examples
		{"add_mul", "3+2*2", 7},
		{"div_only", "3/2", 1},
		{"add_div", "3+5/2", 5}, // 3 + 2 = 5 (5/2=2 integer division)

		// Operator precedence
		{"mul_before_add", "2+3*4", 14},
		{"div_before_sub", "10-6/3", 8},
		{"mul_and_div", "6*2/4", 3},
		{"all_four_ops", "14-3/2*4+5", 15}, // 14 - (3/2)*4 + 5 = 14 - 1*4 + 5 = 15

		// Single operations
		{"single_add", "1+2", 3},
		{"single_sub", "5-3", 2},
		{"single_mul", "4*3", 12},
		{"single_div", "10/3", 3},
		{"single_number", "42", 42},

		// Spaces
		{"spaces_around_ops", " 3 / 2 ", 1},
		{"leading_trailing_spaces", " 3+5 / 2 ", 5},

		// Chain same operation
		{"chain_add", "1+2+3+4", 10},
		{"chain_mul", "2*3*4", 24},
		{"chain_div", "100/10/2", 5},

		// Negative results
		{"negative_result", "1-5", -4},
		{"negative_intermediate", "3*2-10", -4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, calculateII(tt.input))
		})
	}
}
