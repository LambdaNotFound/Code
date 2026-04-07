package graph

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// normalizeAccounts sorts each account's emails (keeping the name first) and
// sorts the account list for deterministic comparison.
func normalizeAccounts(accounts [][]string) [][]string {
	for i := range accounts {
		sort.Strings(accounts[i][1:])
	}
	sort.Slice(accounts, func(i, j int) bool {
		if accounts[i][0] != accounts[j][0] {
			return accounts[i][0] < accounts[j][0]
		}
		if len(accounts[i]) > 1 && len(accounts[j]) > 1 {
			return accounts[i][1] < accounts[j][1]
		}
		return len(accounts[i]) < len(accounts[j])
	})
	return accounts
}

func Test_accountsMerge(t *testing.T) {
	testCases := []struct {
		name     string
		accounts [][]string
		expected [][]string
	}{
		{
			name: "leetcode_example1",
			accounts: [][]string{
				{"John", "johnsmith@mail.com", "john00@mail.com"},
				{"John", "johnnybravo@mail.com"},
				{"John", "johnsmith@mail.com", "john_newyork@mail.com"},
				{"Mary", "mary@mail.com"},
			},
			expected: [][]string{
				{"John", "john00@mail.com", "john_newyork@mail.com", "johnsmith@mail.com"},
				{"John", "johnnybravo@mail.com"},
				{"Mary", "mary@mail.com"},
			},
		},
		{
			name: "leetcode_example2_disjoint",
			accounts: [][]string{
				{"Gabe", "Gabe0@m.co", "Gabe3@m.co", "Gabe1@m.co"},
				{"Kevin", "Kevin3@m.co", "Kevin5@m.co", "Kevin0@m.co"},
				{"Ethan", "Ethan5@m.co", "Ethan4@m.co", "Ethan0@m.co"},
			},
			expected: [][]string{
				{"Ethan", "Ethan0@m.co", "Ethan4@m.co", "Ethan5@m.co"},
				{"Gabe", "Gabe0@m.co", "Gabe1@m.co", "Gabe3@m.co"},
				{"Kevin", "Kevin0@m.co", "Kevin3@m.co", "Kevin5@m.co"},
			},
		},
		{
			name: "two_accounts_fully_merge",
			accounts: [][]string{
				{"Alex", "a@mail.com", "b@mail.com"},
				{"Alex", "b@mail.com", "c@mail.com"},
			},
			expected: [][]string{
				{"Alex", "a@mail.com", "b@mail.com", "c@mail.com"},
			},
		},
		{
			name: "multiple_disjoint_accounts",
			accounts: [][]string{
				{"A", "a1@mail.com"},
				{"B", "b1@mail.com"},
				{"C", "c1@mail.com"},
			},
			expected: [][]string{
				{"A", "a1@mail.com"},
				{"B", "b1@mail.com"},
				{"C", "c1@mail.com"},
			},
		},
		{
			name: "chain_merge",
			accounts: [][]string{
				{"Tom", "t1@mail.com", "t2@mail.com"},
				{"Tom", "t2@mail.com", "t3@mail.com"},
				{"Tom", "t3@mail.com", "t4@mail.com"},
			},
			expected: [][]string{
				{"Tom", "t1@mail.com", "t2@mail.com", "t3@mail.com", "t4@mail.com"},
			},
		},
		{
			name: "same_name_different_emails_no_merge",
			accounts: [][]string{
				{"Sam", "s1@mail.com"},
				{"Sam", "s2@mail.com"},
			},
			expected: [][]string{
				{"Sam", "s1@mail.com"},
				{"Sam", "s2@mail.com"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeAccounts(accountsMerge(tc.accounts))
			want := normalizeAccounts(tc.expected)
			assert.Equal(t, want, got)
		})
	}
}

func Test_accountsMergeUF(t *testing.T) {
	testCases := []struct {
		name     string
		accounts [][]string
		expected [][]string
	}{
		{
			name: "leetcode_example1",
			accounts: [][]string{
				{"John", "johnsmith@mail.com", "john00@mail.com"},
				{"John", "johnnybravo@mail.com"},
				{"John", "johnsmith@mail.com", "john_newyork@mail.com"},
				{"Mary", "mary@mail.com"},
			},
			expected: [][]string{
				{"John", "john00@mail.com", "john_newyork@mail.com", "johnsmith@mail.com"},
				{"John", "johnnybravo@mail.com"},
				{"Mary", "mary@mail.com"},
			},
		},
		{
			name: "two_accounts_fully_merge",
			accounts: [][]string{
				{"Alex", "a@mail.com", "b@mail.com"},
				{"Alex", "b@mail.com", "c@mail.com"},
			},
			expected: [][]string{
				{"Alex", "a@mail.com", "b@mail.com", "c@mail.com"},
			},
		},
		{
			name: "multiple_disjoint_accounts",
			accounts: [][]string{
				{"A", "a1@mail.com"},
				{"B", "b1@mail.com"},
				{"C", "c1@mail.com"},
			},
			expected: [][]string{
				{"A", "a1@mail.com"},
				{"B", "b1@mail.com"},
				{"C", "c1@mail.com"},
			},
		},
		{
			name: "chain_merge",
			accounts: [][]string{
				{"Tom", "t1@mail.com", "t2@mail.com"},
				{"Tom", "t2@mail.com", "t3@mail.com"},
				{"Tom", "t3@mail.com", "t4@mail.com"},
			},
			expected: [][]string{
				{"Tom", "t1@mail.com", "t2@mail.com", "t3@mail.com", "t4@mail.com"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeAccounts(accountsMergeUF(tc.accounts))
			want := normalizeAccounts(tc.expected)
			assert.Equal(t, want, got)
		})
	}
}

func Test_UnionFind(t *testing.T) {
	t.Run("find_returns_own_index", func(t *testing.T) {
		uf := NewUnionFind(3)
		assert.Equal(t, 0, uf.Find(0))
		assert.Equal(t, 1, uf.Find(1))
		assert.Equal(t, 2, uf.Find(2))
	})

	t.Run("union_merges_two_sets", func(t *testing.T) {
		uf := NewUnionFind(3)
		uf.Union(0, 1)
		assert.Equal(t, uf.Find(0), uf.Find(1))
		assert.NotEqual(t, uf.Find(0), uf.Find(2))
	})

	t.Run("union_is_idempotent", func(t *testing.T) {
		uf := NewUnionFind(2)
		uf.Union(0, 1)
		root := uf.Find(0)
		uf.Union(0, 1)
		assert.Equal(t, root, uf.Find(0))
	})

	t.Run("path_compression_flattens_chain", func(t *testing.T) {
		uf := NewUnionFind(4)
		uf.Union(0, 1)
		uf.Union(1, 2)
		uf.Union(2, 3)
		root := uf.Find(3)
		assert.Equal(t, root, uf.Find(0))
		assert.Equal(t, root, uf.Find(1))
		assert.Equal(t, root, uf.Find(2))
	})
}
