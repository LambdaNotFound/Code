package interview

import (
	"sort"
	"sync"
)

/*
Given a list of shopping sessions, where each session is a list of stores
a customer visited, find for each store the list of other stores that
co-occur with it across sessions.

Input:

	records = [
	  ["A", "B", "C"],
	  ["A", "B"],
	  ["B", "C", "D"],
	  ["A", "C"],
	]
	threshold = 2  # minimum co-occurrence count to be "correlated"

Output (for each store, sorted list of correlated stores):

	A: [B, C]     # A+B appears 2x, A+C appears 2x
	B: [A, C]     # B+A appears 2x, B+C appears 2x
	C: [A, B]     # C+A appears 2x, C+B appears 2x
	D: []         # D+anything appears only 1x

follow-up:
1. Top K correlated stores per store instead of threshold <- per store
2. Weighted by session recency — more recent sessions count more: pairWeight[[2]string{a, b}] += weight  // accumulate weighted score
3. Scale — what if there are millions of sessions? (partitioned pair counting, like MapReduce word count but for pairs)
*/

// Time: O(S × K²) where S = sessions, K = avg session size
// Space: O(P) where P = unique pairs
func highCorrelationSellers(records [][]string, threshold int) map[string][]string {
	// count co-occurrences for each pair
	pairCount := map[[2]string]int{}
	storeSet := map[string]bool{}

	for _, session := range records {
		for _, store := range session {
			storeSet[store] = true
		}
		// count every pair in this session
		for i := 0; i < len(session); i++ {
			for j := i + 1; j < len(session); j++ {
				a, b := session[i], session[j]
				if a > b {
					a, b = b, a // normalize order
				}
				pairCount[[2]string{a, b}]++
			}
		}
	}

	// build result: for each store, collect correlated stores
	result := map[string][]string{}
	for store := range storeSet {
		result[store] = []string{}
	}

	for pair, count := range pairCount {
		if count >= threshold {
			a, b := pair[0], pair[1]
			result[a] = append(result[a], b)
			result[b] = append(result[b], a)
		}
	}

	// sort each store's correlated list
	for store := range result {
		sort.Strings(result[store])
	}

	return result
}

func topKCorrelatedSellers(records [][]string, k int) map[string][]string {
	pairCount := map[[2]string]int{}
	storeSet := map[string]bool{}

	for _, session := range records {
		for _, store := range session {
			storeSet[store] = true
		}
		for i := 0; i < len(session); i++ {
			for j := i + 1; j < len(session); j++ {
				a, b := session[i], session[j]
				if a > b {
					a, b = b, a
				}
				pairCount[[2]string{a, b}]++
			}
		}
	}

	// build adjacency: store → list of (neighbor, count)
	type neighbor struct {
		store string
		count int
	}
	adj := map[string][]neighbor{}
	for store := range storeSet {
		adj[store] = []neighbor{}
	}
	for pair, count := range pairCount {
		a, b := pair[0], pair[1]
		adj[a] = append(adj[a], neighbor{b, count})
		adj[b] = append(adj[b], neighbor{a, count})
	}

	// for each store, heap-select top K by count
	result := map[string][]string{}
	for store, neighbors := range adj {
		// If K is large and you want O(n log k) instead of O(n log n),
		// swap sort.Slice for a min-heap of size K
		sort.Slice(neighbors, func(i, j int) bool { // sort descending by count, break ties by name
			if neighbors[i].count != neighbors[j].count {
				return neighbors[i].count > neighbors[j].count
			}
			return neighbors[i].store < neighbors[j].store
		})
		top := []string{}
		for i := 0; i < k && i < len(neighbors); i++ {
			top = append(top, neighbors[i].store)
		}
		result[store] = top
	}

	return result
}

// MapReduce
// Sessions too large to store: Stream from Kafka, process in micro-batches
// Historical + recent blend: Lambda architecture: batch job for historical, stream for recent, merge at read time
func scaledCorrelation(sessionsCh <-chan []string, threshold int, numWorkers int) map[string][]string {
	// each worker processes a chunk of sessions
	partialCounts := make([]map[[2]string]int, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			local := map[[2]string]int{}
			for session := range sessionsCh {
				for i := 0; i < len(session); i++ {
					for j := i + 1; j < len(session); j++ {
						a, b := session[i], session[j]
						if a > b {
							a, b = b, a
						}
						local[[2]string{a, b}]++
					}
				}
			}
			partialCounts[workerID] = local
		}(i)
	}

	wg.Wait()

	// merge partial counts
	merged := map[[2]string]int{}
	for _, partial := range partialCounts {
		for pair, count := range partial {
			merged[pair] += count
		}
	}

	// apply threshold
	result := map[string][]string{}
	for pair, count := range merged {
		if count >= threshold {
			a, b := pair[0], pair[1]
			result[a] = append(result[a], b)
			result[b] = append(result[b], a)
		}
	}
	return result
}
