package backtracking

import "sort"

/**
 * 1087. Brace Expansion
 *
 * If s = "a{b,c}", the first position must be 'a',
 *     but the second position can be either 'b' or 'c'. This produces the words ["ab", "ac"]
 * If s = "{a,b}c{d,e}", you can form: ["acd", "ace", "bcd", "bce"]
 *
 */

/*
Time complexity: O(P × k × log P)

Phase	Cost	Reason
Parsing	O(|s|)	Single pass
Backtracking	O(P × k)	P words, each built character-by-character over k positions
Sorting	O(P × k × log P)	P comparisons, each string comparison is O(k)
Sorting dominates.

Space complexity: O(P × k)

candidates: O(|s|) — the flat list of options
bytes buffer: O(k) — reused in-place via the closure (nice)
res output: O(P × k) — P strings each of length k
*/
func expand(s string) []string {
	candidates := [][]byte{}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		list := []byte{}
		if ch == '{' {
			i++
			for s[i] != '}' {
				if s[i] != ',' {
					list = append(list, s[i])
				}
				i++
			}
		} else {
			list = append(list, ch)
		}
		candidates = append(candidates, list)
	}

	res := []string{}
	bytes := []byte{}
	var backtracking func(pos int)
	backtracking = func(pos int) {
		if pos == len(candidates) {
			res = append(res, string(bytes))
			return
		}

		list := candidates[pos]
		for i := 0; i < len(list); i++ {
			bytes = append(bytes, list[i])
			backtracking(pos + 1)
			bytes = bytes[:len(bytes)-1]
		}
	}

	backtracking(0)
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res
}
