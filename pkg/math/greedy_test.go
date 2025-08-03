package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_leastInterval(t *testing.T) {
    testCases := []struct {
        name     string
        tasks    []byte
        n        int
        expected int
    }{
        {
            "case 1",
            []byte{'A','A','A','B','B','B'},
            2,
            8,
        },
        {
            "case 2",
            []byte{'A','C','A','B','D','B'},
            1,
            6,
        },
        {
            "case 3",
            []byte{'A','A','A', 'B','B','B'},
            3,
            10,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := leastInterval(tc.tasks, tc.n)
            assert.Equal(t, tc.expected, result)
        })
    }
}
