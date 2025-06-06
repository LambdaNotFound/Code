package graph

import "sort"

/**
 * 721. Accounts Merge
 */
 type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(n int) *UnionFind {
    parent, rank := make([]int, n), make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
    }
    return &UnionFind{
        parent: parent,
        rank:   rank,
    }
}

func (uf *UnionFind) Find(node int) int {
    if uf.parent[node] != node {
        uf.parent[node] = uf.Find(uf.parent[node])
    }
    return uf.parent[node]
}

func (uf *UnionFind) Union(p, q int) bool {
    rootP, rootQ := uf.Find(p), uf.Find(q)
    if rootP == rootQ {
        return false
    }
    if uf.rank[rootP] < uf.rank[rootQ] {
        uf.parent[rootP] = rootQ
    } else if uf.rank[rootP] > uf.rank[rootQ] {
        uf.parent[rootQ] = rootP
    } else {
        uf.parent[rootQ] = rootP
        uf.rank[rootP]++
    }

    return true
}

func accountsMerge(accounts [][]string) [][]string {
    uf := NewUnionFind(len(accounts))
    emailToAccount := make(map[string]int)
    for i, account := range accounts {
        for _, email := range account[1:] {
            if acctID, exists := emailToAccount[email]; exists {
                uf.Union(i, acctID)
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
    ret := make([][]string, 0, len(emailGroups))
    for acctID, emails := range emailGroups {
        sort.Strings(emails)
        account := append([]string{accounts[acctID][0]}, emails...)
        ret = append(ret, account)
    }
    return ret
}
