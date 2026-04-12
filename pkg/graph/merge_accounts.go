package graph

import "sort"

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
type UnionFind struct {
	parent []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent}
}

func (uf *UnionFind) Find(node int) int {
	if uf.parent[node] != node {
		uf.parent[node] = uf.Find(uf.parent[node]) // path compression
	}
	return uf.parent[node]
}

func (uf *UnionFind) Union(p, q int) {
	rootP, rootQ := uf.Find(p), uf.Find(q)
	if rootP == rootQ {
		return
	}
	uf.parent[rootP] = rootQ
}

func accountsMergeUF(accounts [][]string) [][]string {
	uf := NewUnionFind(len(accounts))

	emailToAccount := make(map[string]int)
	for i, account := range accounts {
		for _, email := range account[1:] {
			if ownerID, exists := emailToAccount[email]; exists {
				uf.Union(i, ownerID)
			} else {
				emailToAccount[email] = i
			}
		}
	}

	emailGroups := make(map[int][]string)
	for email, accountID := range emailToAccount {
		root := uf.Find(accountID)
		emailGroups[root] = append(emailGroups[root], email)
	}

	result := make([][]string, 0, len(emailGroups))
	for rootID, emails := range emailGroups {
		sort.Strings(emails)
		merged := append([]string{accounts[rootID][0]}, emails...)
		result = append(result, merged)
	}

	return result
}
