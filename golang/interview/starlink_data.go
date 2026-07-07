package interview

import (
	"fmt"
	"sort"
)

type Window struct {
	start, end, data int
	satID            int
}

// lastNonOverlapping returns the largest index j < i whose window ends at or
// before windows[i] starts, or -1 if none. Assumes sorted by end ascending.
func lastNonOverlapping(windows []Window, i int) int {
	lo, hi, res := 0, i-1, -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if windows[mid].end <= windows[i].start {
			res = mid // candidate — look for a later one
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return res
}

// maxData returns the max total data and the chosen windows, given the station
// serves one satellite at a time (chosen windows must be globally non-overlapping).
func maxData(windows []Window) (int, []Window) {
	n := len(windows)
	if n == 0 {
		return 0, nil
	}

	// 1. sort by end time — the foundation of interval scheduling
	sort.Slice(windows, func(i, j int) bool { return windows[i].end < windows[j].end })

	// 2. precompute the last non-overlapping window for each i
	p := make([]int, n)
	for i := range windows {
		p[i] = lastNonOverlapping(windows, i)
	}

	// 3. dp[i] = best total data using windows[0..i]
	dp := make([]int, n)
	dp[0] = windows[0].data
	for i := 1; i < n; i++ {
		include := windows[i].data
		if p[i] >= 0 {
			include += dp[p[i]]
		}
		exclude := dp[i-1]
		if include > exclude {
			dp[i] = include
		} else {
			dp[i] = exclude
		}
	}

	// 4. backtrack to recover the chosen set
	var chosen []Window
	for i := n - 1; i >= 0; {
		include := windows[i].data
		if p[i] >= 0 {
			include += dp[p[i]]
		}
		exclude := 0
		if i > 0 {
			exclude = dp[i-1]
		}
		if include > exclude {
			chosen = append(chosen, windows[i])
			i = p[i] // jump past everything that overlaps
		} else {
			i--
		}
	}

	// chosen came out latest-first — reverse to chronological order
	for l, r := 0, len(chosen)-1; l < r; l, r = l+1, r-1 {
		chosen[l], chosen[r] = chosen[r], chosen[l]
	}
	return dp[n-1], chosen
}

func maxDataTest() {
	windows := []Window{
		{start: 0, end: 3, data: 10, satID: 1},
		{start: 1, end: 4, data: 15, satID: 2},
		{start: 2, end: 5, data: 8, satID: 3},
		{start: 5, end: 8, data: 20, satID: 1},
		{start: 6, end: 9, data: 25, satID: 2},
		{start: 7, end: 10, data: 30, satID: 3},
	}

	total, chosen := maxData(windows)
	fmt.Printf("Max data: %d GB\n", total)
	fmt.Println("Chosen windows:")
	for _, w := range chosen {
		fmt.Printf("  Sat %d: [%d-%d] %d GB\n", w.satID, w.start, w.end, w.data)
	}
}
