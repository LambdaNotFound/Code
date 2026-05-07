package tree

import (
	. "gocode/types"
	"sort"
)

/**
 * 987. Vertical Order Traversal of a Binary Tree
 */
func verticalTraversal(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}

	type pair struct {
		val, row, col int
	}
	pairs := []pair{}

	queue := []*TreeNode{root}
	pairMap := map[*TreeNode]pair{root: pair{root.Val, 0, 0}}
	columns := map[int][]int{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		r, c := pairMap[cur].row, pairMap[cur].col
		pairs = append(pairs, pair{cur.Val, r, c})
		columns[c] = append(columns[c], cur.Val)

		if cur.Left != nil {
			pairMap[cur.Left] = pair{cur.Left.Val, r + 1, c - 1}
			queue = append(queue, cur.Left)
		}
		if cur.Right != nil {
			pairMap[cur.Right] = pair{cur.Right.Val, r + 1, c + 1}
			queue = append(queue, cur.Right)
		}
	}

	columnKeys := make([]int, 0, len(columns))
	for key := range columns {
		columnKeys = append(columnKeys, key)
	}
	sort.Ints(columnKeys)

	res := make([][]int, len(columnKeys))
	for i, k := range columnKeys {
		res[i] = columns[k]
	}
	return res
}
