package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mostVisitedPattern(t *testing.T) {
	tests := []struct {
		name      string
		username  []string
		timestamp []int
		website   []string
		expected  []string
	}{
		{
			// LeetCode example: joe visits [home, about, career]; james visits [home, cart, maps, home, about]
			// joe's pattern (home,about,career)=1; james's patterns include (home,cart,maps),(home,cart,home)... etc.
			// (home,about,career) appears only for joe; pick lexicographically smallest with max count
			name:      "leetcode_example",
			username:  []string{"joe", "joe", "joe", "james", "james", "james", "james", "mary", "mary", "mary"},
			timestamp: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			website:   []string{"home", "about", "career", "home", "cart", "maps", "home", "home", "about", "career"},
			expected:  []string{"home", "about", "career"},
		},
		{
			// Two users share one 3-sequence pattern → count=2 wins
			name:      "shared_pattern_wins",
			username:  []string{"alice", "alice", "alice", "bob", "bob", "bob"},
			timestamp: []int{1, 2, 3, 4, 5, 6},
			website:   []string{"a", "b", "c", "a", "b", "c"},
			expected:  []string{"a", "b", "c"},
		},
		{
			// Single user with exactly 3 visits
			name:      "single_user_3_visits",
			username:  []string{"u", "u", "u"},
			timestamp: []int{3, 1, 2},
			website:   []string{"z", "x", "y"},
			expected:  []string{"x", "y", "z"},
		},
		{
			// Lexicographic tiebreak: two patterns tied at count=1, pick lex smallest
			name:      "lexicographic_tiebreak",
			username:  []string{"alice", "alice", "alice", "bob", "bob", "bob"},
			timestamp: []int{1, 2, 3, 4, 5, 6},
			website:   []string{"a", "b", "c", "d", "e", "f"},
			expected:  []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, mostVisitedPattern(tt.username, tt.timestamp, tt.website))
			assert.Equal(t, tt.expected, mostVisitedPatternClaude(tt.username, tt.timestamp, tt.website))
		})
	}
}
