package utils

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type String string

func foo(str string) string {
	for i := 0; i < len(str); i++ {
		b := str[i] // byte (uint8) everything’s ASCII, so 1 byte = 1 rune.
		fmt.Printf("s[%d] = %c (byte value: %d)\n", i, b, b)
	}

	for i, r := range str {
		fmt.Printf("index %d: rune %c (Unicode: %U)\n", i, r, r)
	}

	return string(str) + " method on custom type"
}

func substr(str string, start, end int) string {
	runes := []rune(str) // convert to rune slice
	return string(runes[start:end])
}

func versionSort(strs []string) {
	sort.Slice(strs, func(i, j int) bool {
		// find where digits start in each string
		splitIdx := func(s string) int {
			for k := 0; k < len(s); k++ {
				if s[k] >= '0' && s[k] <= '9' {
					return k
				}
			}
			return len(s) // no digits
		}

		a, b := strs[i], strs[j]
		ai, bi := splitIdx(a), splitIdx(b)

		// compare alpha prefix first
		if a[:ai] != b[:bi] {
			return a[:ai] < b[:bi]
		}

		// same prefix → compare numeric suffix as integers
		numA, _ := strconv.Atoi(a[ai:])
		numB, _ := strconv.Atoi(b[bi:])
		return numA < numB
	})
}

func Test_string_rune(t *testing.T) {
	foo("你好世界")
}

func Test_substr(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		start    int
		end      int
		expected string
	}{
		{"ascii_full", "hello", 0, 5, "hello"},
		{"ascii_partial", "hello", 1, 4, "ell"},
		{"unicode_full", "你好世界", 0, 4, "你好世界"},
		{"unicode_partial", "你好世界", 1, 3, "好世"},
		{"empty_result", "hello", 2, 2, ""},
		{"single_rune", "你好世界", 0, 1, "你"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, substr(tc.str, tc.start, tc.end))
		})
	}
}
