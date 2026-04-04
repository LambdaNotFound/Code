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

func Test_calculateIII(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// LeetCode 772 examples
		{"leetcode_1", "1+1", 2},
		{"leetcode_2", "6-4/2", 4},   // 6 - 2 = 4
		{"leetcode_3", "2*(5+5*2)/3+(6/2+8)", 21},
		{"leetcode_4", "(2+6*3+5-(3*14/7+2)*5)+3", -12},
		{"leetcode_5", "2*(3+4)-1", 13},

		// Operator precedence (* and / before + and -)
		{"mul_before_add", "2+3*4", 14},
		{"div_before_sub", "10-6/3", 8},
		{"chain_mul_div", "6*2/4", 3},    // left-to-right: (6*2)/4 = 3
		{"chain_div_mul", "8/2*3", 12},   // left-to-right: (8/2)*3 = 12

		// Parentheses override precedence
		{"parens_add_first", "(1+2)*3", 9},
		{"parens_vs_mul", "(2+3)*(4+5)", 45},
		{"parens_in_div", "10/(2+3)", 2},     // 10/5 = 2
		{"deep_parens", "((1+2))", 3},
		{"nested_mul", "((3+2)*2)", 10},

		// All four operators with parentheses
		{"mixed_all", "(1+2)*(3+(4*5))", 69}, // 3*(3+20) = 3*23 = 69
		{"sub_in_parens", "(10-3)*2", 14},
		{"div_in_parens", "(10/2)+3", 8},

		// Integer division truncates toward zero
		{"int_div", "7/2", 3},
		{"int_div_exact", "10/2", 5},
		{"int_div_chain", "100/10/2", 5},

		// Single operations and numbers
		{"single_number", "5", 5},
		{"single_add", "1+2", 3},
		{"single_mul", "4*3", 12},
		{"single_div", "10/3", 3},

		// Negative results
		{"negative_result", "1-5", -4},
		{"negative_mul", "10-3*4", -2}, // 10 - 12 = -2
		{"zero_result", "3-3", 0},

		// Multi-digit numbers
		{"multi_digit", "100*2+50/5", 210}, // 200 + 10 = 210
		{"large_parens", "(100+200)*3", 900},

		// Chain same operation
		{"chain_add", "1+2+3+4", 10},
		{"chain_mul", "1*2*3*4", 24},
		{"chain_sub", "10-1-2-3", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, calculateIII(tt.input), "calculateIII")
			assert.Equal(t, tt.expected, calculateClaude(tt.input), "calculateClaude")
		})
	}
}
