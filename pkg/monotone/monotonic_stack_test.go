package monotone

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_trap(t *testing.T) {
    testCases := []struct {
        name     string
        height   []int
        expected int
    }{
        {
            "case 1",
            []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1},
            6,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := trap(tc.height)
            assert.Equal(t, tc.expected, result)
        })
    }
}
