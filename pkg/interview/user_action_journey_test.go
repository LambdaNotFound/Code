package interview

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_groupLogsByUser(t *testing.T) {
	testCases := []struct {
		name     string
		logs     []log
		expected map[int][]log // userId → sorted logs
	}{
		{
			name: "single user logs sorted by time",
			logs: []log{
				{100, 1300, "C"},
				{100, 1000, "A"},
				{100, 1200, "B"},
			},
			expected: map[int][]log{
				100: {{100, 1000, "A"}, {100, 1200, "B"}, {100, 1300, "C"}},
			},
		},
		{
			name: "multiple users grouped and sorted independently",
			logs: []log{
				{100, 1000, "A"},
				{200, 1100, "A"},
				{200, 1200, "B"},
				{100, 1200, "B"},
				{100, 1300, "C"},
				{200, 1400, "A"},
			},
			expected: map[int][]log{
				100: {{100, 1000, "A"}, {100, 1200, "B"}, {100, 1300, "C"}},
				200: {{200, 1100, "A"}, {200, 1200, "B"}, {200, 1400, "A"}},
			},
		},
		{
			name: "single log entry",
			logs: []log{
				{300, 1500, "B"},
			},
			expected: map[int][]log{
				300: {{300, 1500, "B"}},
			},
		},
		{
			name: "same user same-time logs preserved",
			logs: []log{
				{100, 1000, "A"},
				{100, 1000, "B"},
			},
			expected: map[int][]log{
				100: {{100, 1000, "A"}, {100, 1000, "B"}},
			},
		},
		{
			name:     "empty logs returns empty groups",
			logs:     []log{},
			expected: map[int][]log{},
		},
		{
			name: "three users grouped independently",
			logs: []log{
				{100, 1000, "A"},
				{200, 1100, "B"},
				{300, 1200, "C"},
				{100, 1300, "B"},
				{300, 1400, "A"},
			},
			expected: map[int][]log{
				100: {{100, 1000, "A"}, {100, 1300, "B"}},
				200: {{200, 1100, "B"}},
				300: {{300, 1200, "C"}, {300, 1400, "A"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := groupLogsByUser(tc.logs)

			// outer order is non-deterministic (map iteration); index by userId
			grouped := make(map[int][]log)
			for _, group := range result {
				grouped[group[0].userId] = group
			}

			assert.Equal(t, len(tc.expected), len(grouped))
			for userId, expectedLogs := range tc.expected {
				assert.Equal(t, expectedLogs, grouped[userId])
			}
		})
	}
}

func Test_Trie_Insert(t *testing.T) {
	testCases := []struct {
		name        string
		logGroups   [][]log
		checkAction string              // top-level action to inspect
		wantCount   int                 // expected count at that node
		wantChildren []string           // expected children of that node
	}{
		{
			name: "single sequence",
			logGroups: [][]log{
				{{100, 1000, "A"}, {100, 1200, "B"}},
			},
			checkAction:  "A",
			wantCount:    1,
			wantChildren: []string{"B"},
		},
		{
			name: "two sequences sharing a prefix increments count",
			logGroups: [][]log{
				{{100, 1000, "A"}, {100, 1200, "B"}},
				{{200, 1100, "A"}, {200, 1200, "B"}, {200, 1400, "A"}},
			},
			checkAction:  "A",
			wantCount:    2,
			wantChildren: []string{"B"},
		},
		{
			name: "diverging sequences create sibling children",
			logGroups: [][]log{
				{{100, 1000, "A"}, {100, 1200, "B"}},
				{{200, 1100, "A"}, {200, 1200, "C"}},
			},
			checkAction:  "A",
			wantCount:    2,
			wantChildren: []string{"B", "C"},
		},
		{
			name: "three-action sequence builds deep path",
			logGroups: [][]log{
				{{100, 1000, "A"}, {100, 1200, "B"}, {100, 1300, "C"}},
			},
			checkAction:  "A",
			wantCount:    1,
			wantChildren: []string{"B"},
		},
		{
			name: "repeated action in sequence stored as self-child",
			logGroups: [][]log{
				{{300, 1500, "B"}, {300, 1550, "B"}},
			},
			checkAction:  "B",
			wantCount:    1,
			wantChildren: []string{"B"}, // B → B
		},
		{
			name: "count accumulates for three sequences sharing root action",
			logGroups: [][]log{
				{{100, 1000, "A"}, {100, 1200, "B"}},
				{{200, 1100, "A"}, {200, 1200, "B"}},
				{{300, 1050, "A"}, {300, 1200, "C"}},
			},
			checkAction:  "A",
			wantCount:    3,
			wantChildren: []string{"B", "C"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			root := &Trie{children: make(map[string]*Trie)}
			for _, group := range tc.logGroups {
				root.Insert(group)
			}

			node, found := root.children[tc.checkAction]
			assert.True(t, found, "expected action %q at root", tc.checkAction)
			assert.Equal(t, tc.wantCount, node.count)

			childKeys := make([]string, 0, len(node.children))
			for k := range node.children {
				childKeys = append(childKeys, k)
			}
			sort.Strings(childKeys)
			sort.Strings(tc.wantChildren)
			assert.Equal(t, tc.wantChildren, childKeys)
		})
	}
}
