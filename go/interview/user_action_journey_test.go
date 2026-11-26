package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseLog(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected []LogEntry
	}{
		{
			name: "single_line",
			raw:  "100 1000 A",
			expected: []LogEntry{
				{user: 100, time: 1000, action: "A"},
			},
		},
		{
			name: "multi_line",
			raw: `100 1000 A
200 1100 B
300 1200 C`,
			expected: []LogEntry{
				{user: 100, time: 1000, action: "A"},
				{user: 200, time: 1100, action: "B"},
				{user: 300, time: 1200, action: "C"},
			},
		},
		{
			name: "sample_from_comments",
			raw: `100  1000 A
200  1100 A
200  1200 B
100  1200 B
100  1300 C
200  1400 A
300  1500 B
300  1550 B`,
			expected: []LogEntry{
				{user: 100, time: 1000, action: "A"},
				{user: 200, time: 1100, action: "A"},
				{user: 200, time: 1200, action: "B"},
				{user: 100, time: 1200, action: "B"},
				{user: 100, time: 1300, action: "C"},
				{user: 200, time: 1400, action: "A"},
				{user: 300, time: 1500, action: "B"},
				{user: 300, time: 1550, action: "B"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, parseLog(tt.raw))
		})
	}
}

func Test_buildUserSequences(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		result := buildUserSequences([]LogEntry{})
		assert.Empty(t, result)
	})

	t.Run("single_user_already_sorted", func(t *testing.T) {
		entries := []LogEntry{
			{user: 1, time: 100, action: "A"},
			{user: 1, time: 200, action: "B"},
			{user: 1, time: 300, action: "C"},
		}
		result := buildUserSequences(entries)
		assert.Equal(t, map[int][]ActionEntry{
			1: {{"A", 100}, {"B", 200}, {"C", 300}},
		}, result)
	})

	t.Run("single_user_needs_sort", func(t *testing.T) {
		entries := []LogEntry{
			{user: 1, time: 300, action: "C"},
			{user: 1, time: 100, action: "A"},
			{user: 1, time: 200, action: "B"},
		}
		result := buildUserSequences(entries)
		assert.Equal(t, map[int][]ActionEntry{
			1: {{"A", 100}, {"B", 200}, {"C", 300}},
		}, result)
	})

	t.Run("multi_user", func(t *testing.T) {
		entries := []LogEntry{
			{user: 1, time: 200, action: "B"},
			{user: 2, time: 500, action: "X"},
			{user: 1, time: 100, action: "A"},
			{user: 2, time: 400, action: "Y"},
		}
		result := buildUserSequences(entries)
		assert.Equal(t, []ActionEntry{{"A", 100}, {"B", 200}}, result[1])
		assert.Equal(t, []ActionEntry{{"Y", 400}, {"X", 500}}, result[2])
	})
}

func Test_trieInsert(t *testing.T) {
	t.Run("single_sequence", func(t *testing.T) {
		root := newTrieNode()
		root.insert([]ActionEntry{{"A", 1000}, {"B", 1200}, {"C", 1300}})

		assert.Equal(t, 1, root.children["A"].count)
		assert.Equal(t, 1000, root.firstSeen["A"])

		nodeA := root.children["A"]
		assert.Equal(t, 1, nodeA.children["B"].count)
		assert.Equal(t, 1200, nodeA.firstSeen["B"])

		nodeB := nodeA.children["B"]
		assert.Equal(t, 1, nodeB.children["C"].count)
		assert.Equal(t, 1300, nodeB.firstSeen["C"])
	})

	t.Run("two_sequences_shared_prefix", func(t *testing.T) {
		root := newTrieNode()
		root.insert([]ActionEntry{{"A", 1000}, {"B", 1200}})
		root.insert([]ActionEntry{{"A", 1100}, {"B", 1400}})

		// A appears twice at root level
		assert.Equal(t, 2, root.children["A"].count)
		// firstSeen should be the earlier timestamp
		assert.Equal(t, 1000, root.firstSeen["A"])

		nodeA := root.children["A"]
		assert.Equal(t, 2, nodeA.children["B"].count)
		assert.Equal(t, 1200, nodeA.firstSeen["B"])
	})

	t.Run("first_seen_tracks_earliest", func(t *testing.T) {
		root := newTrieNode()
		// Insert later timestamp first
		root.insert([]ActionEntry{{"A", 2000}})
		root.insert([]ActionEntry{{"A", 1000}})

		assert.Equal(t, 1000, root.firstSeen["A"])
	})

	t.Run("sample_data_full_trie", func(t *testing.T) {
		// Sample from file comments:
		// User 100: A(1000) -> B(1200) -> C(1300)
		// User 200: A(1100) -> B(1200) -> A(1400)
		// User 300: B(1500) -> B(1550)
		sequences := map[int][]ActionEntry{
			100: {{"A", 1000}, {"B", 1200}, {"C", 1300}},
			200: {{"A", 1100}, {"B", 1200}, {"A", 1400}},
			300: {{"B", 1500}, {"B", 1550}},
		}

		root := newTrieNode()
		for _, seq := range sequences {
			root.insert(seq)
		}

		// Root level: A appears 2x (users 100, 200), B appears 1x (user 300)
		assert.Equal(t, 2, root.children["A"].count)
		assert.Equal(t, 1, root.children["B"].count)

		// A -> B: both users 100 and 200 do B after A
		nodeA := root.children["A"]
		assert.Equal(t, 2, nodeA.children["B"].count)

		// A -> B -> C: only user 100
		nodeAB := nodeA.children["B"]
		assert.Equal(t, 1, nodeAB.children["C"].count)

		// A -> B -> A: only user 200
		assert.Equal(t, 1, nodeAB.children["A"].count)

		// B -> B: user 300 does B twice
		nodeB := root.children["B"]
		assert.Equal(t, 1, nodeB.children["B"].count)
	})
}

func Test_newTrieNode(t *testing.T) {
	node := newTrieNode()
	assert.NotNil(t, node)
	assert.NotNil(t, node.children)
	assert.NotNil(t, node.firstSeen)
	assert.Equal(t, 0, node.count)
	assert.Empty(t, node.children)
	assert.Empty(t, node.firstSeen)
}

// ── printTrie (coverage) ──────────────────────────────────────────────────────

func Test_printTrie(t *testing.T) {
	t.Run("empty_trie_prints_nothing", func(t *testing.T) {
		// Should not panic on empty root
		printTrie(newTrieNode(), 0)
	})

	t.Run("single_level_trie", func(t *testing.T) {
		root := newTrieNode()
		root.insert([]ActionEntry{{"A", 1000}})
		printTrie(root, 0)
	})

	t.Run("nested_trie_with_indent", func(t *testing.T) {
		root := newTrieNode()
		root.insert([]ActionEntry{{"A", 1000}, {"B", 1200}, {"C", 1300}})
		root.insert([]ActionEntry{{"A", 1100}, {"B", 1200}, {"A", 1400}})
		printTrie(root, 0)
	})

	t.Run("tie_break_by_action_name", func(t *testing.T) {
		// Two actions at the same timestamp — alphabetical tie-break
		root := newTrieNode()
		root.insert([]ActionEntry{{"Z", 1000}})
		root.insert([]ActionEntry{{"A", 1000}})
		printTrie(root, 0)
	})
}

// ── testUserActions (demo function coverage) ──────────────────────────────────

func Test_testUserActions(t *testing.T) {
	// testUserActions() prints to stdout; calling it exercises its branches.
	testUserActions()
}
