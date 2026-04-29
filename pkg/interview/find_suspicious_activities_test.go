package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * findSuspiciousActivities
 *
 * Each activity is encoded as "name,city,action".
 * The function performs a BFS starting from the seed suspicious activities,
 * expanding to any newActivity that is "similar" (differs in at most k fields).
 * The expected output contains only activities from newActivities that are
 * reachable from the seed set.
 */

func Test_findSuspiciousActivities(t *testing.T) {
	testCases := []struct {
		name                 string
		suspiciousActivities [][]string
		newActivities        [][]string
		k                    int
		expected             [][]string
	}{
		{
			name: "example_from_spec",
			suspiciousActivities: [][]string{
				{"Brad", "San Francisco", "withdraw"},
			},
			newActivities: [][]string{
				{"Joe", "Miami", "withdraw"},
				{"John", "San Francisco", "deposit"},
				{"Albert", "London", "withdraw"},
				{"Diana", "London", "withdraw"},
				{"Diana", "San Francisco", "withdraw"},
				{"Joe", "New York", "updateaddress"},
			},
			k: 2,
			expected: [][]string{
				{"Albert", "London", "withdraw"},
				{"Diana", "London", "withdraw"},
				{"Diana", "San Francisco", "withdraw"},
			},
		},
	}

	impls := []struct {
		name string
		fn   func([][]string, [][]string, int) [][]string
	}{
		{"BFS_quadratic", findSuspiciousActivities},
		{"BFS_inverted_index", findSuspiciousActivitiesOpt},
		{"union_find", findSuspiciousActivitiesUF},
	}

	for _, tc := range testCases {
		for _, impl := range impls {
			t.Run(tc.name+"/"+impl.name, func(t *testing.T) {
				result := impl.fn(tc.suspiciousActivities, tc.newActivities, tc.k)
				assert.ElementsMatch(t, tc.expected, result)
			})
		}
	}
}
