package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_minRemoveToMakeValid(t *testing.T) {
    testCases := []struct {
        name     string
        s        string
        expected string
    }{
        {
            name:     "already valid",
            s:        "a(b)c",
            expected: "a(b)c",
        },
        {
            name:     "remove extra closing",
            s:        "a)b(c)d",
            expected: "ab(c)d",
        },
        {
            name:     "remove extra opening",
            s:        "a(b(c",
            expected: "a(bc",
        },
        {
            name:     "mixed invalid",
            s:        ")a(b))c(",
            expected: "a(b)c",
        },
        {
            name:     "only parentheses open",
            s:        "(()",
            expected: "()",
        },
        {
            name:     "only parentheses invalid",
            s:        "))((",
            expected: "",
        },
        {
            name:     "empty string",
            s:        "",
            expected: "",
        },
        {
            name:     "letters only",
            s:        "abcde",
            expected: "abcde",
        },
        {
            name:     "nested complex",
            s:        "a((b)c)d)",
            expected: "a((b)c)d",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := minRemoveToMakeValid(tc.s)
            assert.Equal(t, tc.expected, result)
        })
    }
}

func TestCalculate(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
    }{
        {
            name:     "Simple addition",
            input:    "1 + 1",
            expected: 2,
        },
        {
            name:     "Simple subtraction",
            input:    "2 - 1",
            expected: 1,
        },
        {
            name:     "Mixed + and -",
            input:    "2 - 1 + 2",
            expected: 3,
        },
        {
            name:     "Single number",
            input:    "42",
            expected: 42,
        },
        {
            name:     "Leading and trailing spaces",
            input:    "   3   ",
            expected: 3,
        },
        {
            name:     "Basic parentheses",
            input:    "(1 + 1)",
            expected: 2,
        },
        {
            name:     "Nested parentheses",
            input:    "(1+(4+5+2)-3)+(6+8)",
            expected: 23,
        },
        {
            name:     "Negative result",
            input:    "1 - 5",
            expected: -4,
        },
        {
            name:     "Negative number inside parentheses",
            input:    "1 - (5 - 2)",
            expected: -2,
        },
        {
            name:     "Deep nesting",
            input:    "((((1+2))))",
            expected: 3,
        },
        {
            name:     "Multiple spaces",
            input:    " 2-   1 +   2 ",
            expected: 3,
        },
        {
            name:     "Chain parentheses",
            input:    "(1)-(2)-(3)",
            expected: -4,
        },
        {
            name:     "Unary minus before parentheses",
            input:    "-(3 + 4)",
            expected: -7,
        },
        {
            name:     "Unary minus before number",
            input:    "-3 + 1",
            expected: -2,
        },
        {
            name:     "Unary minus nested",
            input:    "-(1-(2+3))",
            expected: -4,
        },
        {
            name:     "Empty string",
            input:    "",
            expected: 0,
        },
        {
            name:     "Double parentheses",
            input:    "((2))",
            expected: 2,
        },
        {
            name:     "Complex spaces and nesting",
            input:    " ( 7 -( 3 + (2 - 1) ) ) ",
            expected: 3,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got := calculate(tc.input)
            assert.Equal(t, tc.expected, got)
        })
    }
}

func TestEvalRPN(t *testing.T) {
    tests := []struct {
        name   string
        tokens []string
        want   int
    }{
        {
            name:   "basic multiplication",
            tokens: []string{"2", "1", "+", "3", "*"},
            want:   9,
        },
        {
            name:   "division and addition",
            tokens: []string{"4", "13", "5", "/", "+"},
            want:   6,
        },
        {
            name:   "large example",
            tokens: []string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"},
            want:   22,
        },
        {
            name:   "simple subtraction",
            tokens: []string{"3", "4", "-"},
            want:   -1,
        },
        {
            name:   "negative multiplication",
            tokens: []string{"-3", "4", "*"},
            want:   -12,
        },
        {
            name:   "division truncate positive",
            tokens: []string{"4", "3", "/"},
            want:   1,
        },
        {
            name:   "division truncate negative left",
            tokens: []string{"-7", "3", "/"},
            want:   -2,
        },
        {
            name:   "division truncate negative right",
            tokens: []string{"7", "-3", "/"},
            want:   -2,
        },
        {
            name:   "single number",
            tokens: []string{"42"},
            want:   42,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := evalRPN(tt.tokens)
            assert.Equal(t, tt.want, got)
        })
    }
}
