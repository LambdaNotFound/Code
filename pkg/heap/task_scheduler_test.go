package heap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_leastInterval(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []byte
		n        int
		expected int
	}{
		{name: "leetcode_example1", tasks: []byte("AAABBB"), n: 2, expected: 8},
		{name: "leetcode_example2", tasks: []byte("AAABBB"), n: 0, expected: 6},
		{name: "leetcode_example3", tasks: []byte("AAAAAABCDEFG"), n: 2, expected: 16},
		{name: "single_task", tasks: []byte("A"), n: 5, expected: 1},
		{name: "no_cooldown_all_different", tasks: []byte("ABC"), n: 2, expected: 3},
		{name: "triple_a_cooldown_2", tasks: []byte("AAA"), n: 2, expected: 7},
		{name: "all_same_no_cooldown", tasks: []byte("AAA"), n: 0, expected: 3},
		{name: "two_tasks_equal_freq", tasks: []byte("AABB"), n: 2, expected: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := leastInterval(tt.tasks, tt.n)
			assert.Equal(t, tt.expected, result)
		})
	}
}
