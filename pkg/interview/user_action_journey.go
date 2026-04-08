package interview

/*
 Sample log messages:
"""
USER TIME ACTION
100  1000 A
200  1100 A
200  1200 B
100  1200 B
100  1300 C
200  1400 A
300  1500 B
300  1550 B
"""

Sample output:
"""
A (2)
  -> B (2)
     -> C (1)
     -> A (1)
B (1)
  -> B (1)
"""
*/

/**
 * Time Complexity: O(N log N + T × A log A)
 *
 * Parse logsO(N): Single pass
 * Sort each user's logs by time: O(N log N)
 * Across all usersInsert into trie: O(U × L), Each entry traverses L nodes
 * Print trie (sort children at each node): O(T × A log A), T = trie nodes, sort up to A children each
 */

import (
	"fmt"
	"sort"
	"strings"
)

// --- Types ---

type LogEntry struct {
	user, time int
	action     string
}

type ActionEntry struct {
	action string
	time   int
}

type TrieNode struct {
	count     int
	children  map[string]*TrieNode
	firstSeen map[string]int // action -> earliest timestamp at this level
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		children:  map[string]*TrieNode{},
		firstSeen: map[string]int{},
	}
}

// pre-processing
func parseLog(raw string) []LogEntry {
	var entries []LogEntry
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		var user, time int
		var action string
		fmt.Sscanf(strings.TrimSpace(line), "%d %d %s", &user, &time, &action)
		entries = append(entries, LogEntry{user, time, action})
	}
	return entries
}

func buildUserSequences(entries []LogEntry) map[int][]ActionEntry {
	userEntries := map[int][]LogEntry{}
	for _, e := range entries {
		userEntries[e.user] = append(userEntries[e.user], e)
	}
	sequences := map[int][]ActionEntry{}
	for user, logs := range userEntries {
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].time < logs[j].time
		})
		for _, log := range logs {
			sequences[user] = append(sequences[user], ActionEntry{log.action, log.time})
		}
	}
	return sequences
}

// prefix tree ops
func (node *TrieNode) insert(seq []ActionEntry) {
	cur := node
	for _, e := range seq {
		if _, ok := cur.children[e.action]; !ok {
			cur.children[e.action] = newTrieNode()
			cur.firstSeen[e.action] = e.time
		} else if e.time < cur.firstSeen[e.action] {
			cur.firstSeen[e.action] = e.time // always keep earliest
		}
		cur = cur.children[e.action]
		cur.count++
	}
}

func printTrie(node *TrieNode, depth int) {
	type actionTime struct {
		action string
		time   int
	}
	order := make([]actionTime, 0, len(node.firstSeen))
	for action, t := range node.firstSeen {
		order = append(order, actionTime{action, t})
	}
	sort.Slice(order, func(i, j int) bool {
		if order[i].time != order[j].time {
			return order[i].time < order[j].time
		}
		return order[i].action < order[j].action // alphabetical tie-break
	})

	indent := strings.Repeat("   ", depth)
	arrow := ""
	if depth > 0 {
		arrow = "-> "
	}
	for _, at := range order {
		child := node.children[at.action]
		fmt.Printf("%s%s%s (%d)\n", indent, arrow, at.action, child.count)
		printTrie(child, depth+1)
	}
}

func main() {
	raw := `100  1000 A
200  1100 A
200  1200 B
100  1200 B
100  1300 C
200  1400 A
300  1500 B
300  1550 B`

	entries := parseLog(raw)
	sequences := buildUserSequences(entries)
	trie := &TrieNode{children: map[string]*TrieNode{}, firstSeen: map[string]int{}}
	for _, seq := range sequences {
		trie.insert(seq)
	}

	fmt.Println("Journey summary:")
	printTrie(trie, 0)
}
