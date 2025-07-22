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
