package graph

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper: sort each account's emails AND sort the list of accounts
func normalize(accounts [][]string) [][]string {
    for i := range accounts {
        // sort emails except the first element (the name)
        sort.Strings(accounts[i][1:])
    }
    sort.Slice(accounts, func(i, j int) bool {
        if accounts[i][0] == accounts[j][0] {
            // compare first email if names are equal
            return accounts[i][1] < accounts[j][1]
        }
        return accounts[i][0] < accounts[j][0]
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
            name: "leetcode example 1",
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
            name: "leetcode example 2",
            accounts: [][]string{
                {"Gabe", "Gabe0@m.co", "Gabe3@m.co", "Gabe1@m.co"},
                {"Kevin", "Kevin3@m.co", "Kevin5@m.co", "Kevin0@m.co"},
                {"Ethan", "Ethan5@m.co", "Ethan4@m.co", "Ethan0@m.co"},
                {"Hanzo", "Hanzo3@m.co", "Hanzo1@m.co", "Hanzo0@m.co"},
                {"Fern", "Fern5@m.co", "Fern1@m.co", "Fern0@m.co"},
            },
            expected: [][]string{
                {"Ethan", "Ethan0@m.co", "Ethan4@m.co", "Ethan5@m.co"},
                {"Fern", "Fern0@m.co", "Fern1@m.co", "Fern5@m.co"},
                {"Gabe", "Gabe0@m.co", "Gabe1@m.co", "Gabe3@m.co"},
                {"Hanzo", "Hanzo0@m.co", "Hanzo1@m.co", "Hanzo3@m.co"},
                {"Kevin", "Kevin0@m.co", "Kevin3@m.co", "Kevin5@m.co"},
            },
        },
        {
            name: "two accounts fully merge",
            accounts: [][]string{
                {"Alex", "a@mail.com", "b@mail.com"},
                {"Alex", "b@mail.com", "c@mail.com"},
            },
            expected: [][]string{
                {"Alex", "a@mail.com", "b@mail.com", "c@mail.com"},
            },
        },
        {
            name: "multiple disjoint accounts",
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
            name: "chain merge",
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
            name: "same name different accounts (should not merge)",
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
            result := accountsMerge(tc.accounts)

            // normalize order for comparison
            got := normalize(result)
            want := normalize(tc.expected)

            assert.Equal(t, want, got)
        })
    }
}
