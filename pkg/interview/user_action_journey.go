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

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type log struct {
	userId, time int
	action       string
}

type Trie struct {
	action   string
	count    int
	time     int
	children map[string]*Trie
}

func (t *Trie) Insert(logs []log) {
	cur := t
	for _, log := range logs {
		if _, found := cur.children[log.action]; !found {
			new := &Trie{log.action, 0, 0, make(map[string]*Trie)}
			cur.children[log.action] = new
		}
		cur.children[log.action].count++
		cur.children[log.action].time = min(cur.children[log.action].time, log.time)
		cur = cur.children[log.action]
	}
}

func (t *Trie) Print(depth int) {
	children := []*Trie{}
	for _, v := range t.children {
		children = append(children, v)
	}
	sort.Slice(children, func(i, j int) bool { // sort keys
		return children[i].time < children[j].time
	})

	for _, c := range children {
		indent := strings.Repeat("  ", depth)
		if depth == 0 {
			fmt.Printf("%s%v (%d) \n", indent, c.action, c.count)
		} else {
			fmt.Printf("%s-> %v (%d) \n", indent, c.action, c.count)
		}
		c.Print(depth + 1)
	}
}

func groupLogsByUser(logs []log) [][]log {
	hashmap := make(map[int][]log)
	for _, log := range logs {
		hashmap[log.userId] = append(hashmap[log.userId], log)
	}

	groupLogs := [][]log{}
	for _, list := range hashmap {
		sort.Slice(list, func(i, j int) bool {
			return list[i].time < list[j].time
		})
		groupLogs = append(groupLogs, list)
	}

	return groupLogs
}

func main() {
	input := `USER TIME ACTION
	100  1000 A
	200  1100 A
	200  1200 B
	100  1200 B
	100  1300 C
	200  1400 A
	300  1500 B
	300  1550 B`

	rows := strings.Split(input, "\n")
	logs := []log{}
	for _, row := range rows[1:] {
		fields := strings.Fields(row)
		userId, _ := strconv.Atoi(fields[0])
		time, _ := strconv.Atoi(fields[1])
		logs = append(logs, log{userId, time, fields[2]})
	}

	groupLogs := groupLogsByUser(logs)
	for _, l := range groupLogs {
		fmt.Printf("%+v\n", l)
	}
	fmt.Printf("\n")

	root := &Trie{children: make(map[string]*Trie)}
	for _, l := range groupLogs {
		root.Insert(l)
	}
	root.Print(0)
}
