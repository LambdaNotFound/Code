package graph

import (
	"slices"
	"sort"
)

/**
 * 721. Accounts Merge
 *
 * Two accounts definitely belong to the same person if there is some common email to both accounts.
 * Note that even if two accounts have the same name, they may belong to different people as people
 * could have the same name.
 *
 * TimeO(N×K log(N×K)) — N = accounts, K = avg emails, log factor from sorting
 * SpaceO(N×K) — emailToAccounts map
 */
func accountsMerge(accounts [][]string) [][]string {
	// map each email to all account indices that contain it
	emailToAccounts := make(map[string][]int)
	for idx, account := range accounts {
		for _, email := range account[1:] {
			emailToAccounts[email] = append(emailToAccounts[email], idx)
		}
	}

	result := make([][]string, 0)
	visited := make([]bool, len(accounts))

	for idx, account := range accounts {
		if visited[idx] {
			continue
		}

		// BFS to collect all connected accounts
		queue := []int{idx}
		visited[idx] = true
		emails := make(map[string]bool)

		for len(queue) > 0 {
			accountIdx := queue[0]
			queue = queue[1:]

			// traverse all emails of current account
			for _, email := range accounts[accountIdx][1:] {
				if emails[email] {
					continue // skip already collected emails
				}
				emails[email] = true

				// find all accounts sharing this email
				for _, neighborIdx := range emailToAccounts[email] {
					if !visited[neighborIdx] {
						visited[neighborIdx] = true
						queue = append(queue, neighborIdx)
					}
				}
			}
		}

		// build result: sorted emails prepended with account name
		merged := make([]string, 0, len(emails)+1)
		for email := range emails {
			merged = append(merged, email)
		}
		sort.Strings(merged)
		merged = append([]string{account[0]}, merged...)
		result = append(result, merged)
	}

	return result
}

/**
 * Union-Find
 */
func accountsMergeUF(accounts [][]string) [][]string {
	n := len(accounts)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union := func(x, y int) {
		rootX, rootY := find(x), find(y)
		if rootX != rootY {
			parent[rootX] = rootY
		}
	}

	// map each email to all account indices that contain it
	emailToAccounts := make(map[string][]int)
	for i, account := range accounts {
		for _, email := range account[1:] {
			emailToAccounts[email] = append(emailToAccounts[email], i)
		}
	}

	// union all accounts sharing an email
	for _, indices := range emailToAccounts {
		for k := 1; k < len(indices); k++ {
			union(indices[0], indices[k])
		}
	}

	// group emails by root account index
	merged := make(map[int][]string)
	for i, account := range accounts {
		root := find(i)
		merged[root] = append(merged[root], account[1:]...)
	}

	res := [][]string{}
	for root, emails := range merged {
		sort.Strings(emails)
		emails = slices.Compact(emails)
		name := accounts[root][0]
		res = append(res, append([]string{name}, emails...))
	}
	return res
}
