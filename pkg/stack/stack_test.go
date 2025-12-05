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
            "case 1",
            "lee(t(c)o)de)",
            "lee(t(c)o)de",
        },
        {
            "case 2",
            "a)b(c)d",
            "ab(c)d",
        },
        {
            "case 3",
            "))((",
            "",
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
