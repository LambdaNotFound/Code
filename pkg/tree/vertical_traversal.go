package tree

import (
	. "gocode/types"
	"sort"
)

/**
 * 987. Vertical Order Traversal of a Binary Tree
 */

/**
 * Time: O(N log N) Sorting the column keys: O(N log N) worst case
 * Space: O(N + H) map holds all N entries: O(N), DFS call stack: O(H) where H is the tree height
 */
func verticalTraversalDFS(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}

	type entry struct {
		row, val int
	}

	columns := map[int][]entry{}
	var dfs func(node *TreeNode, row, col int)
	dfs = func(node *TreeNode, row, col int) {
		if node == nil {
			return
		}
		columns[col] = append(columns[col], entry{row, node.Val})
		dfs(node.Left, row+1, col-1)
		dfs(node.Right, row+1, col+1)
	}
	dfs(root, 0, 0)

	columnKeys := make([]int, 0, len(columns))
	for key := range columns {
		columnKeys = append(columnKeys, key)
	}
	sort.Ints(columnKeys)

	var res [][]int
	for _, key := range columnKeys {
		entries := columns[key]
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].row != entries[j].row {
				return entries[i].row < entries[j].row
			}
			return entries[i].val < entries[j].val
		})
		vals := make([]int, len(entries))
		for i, e := range entries {
			vals[i] = e.val
		}
		res = append(res, vals)
	}
	return res
}

func verticalTraversalBFS(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}

	type entry struct {
		node          *TreeNode
		col, row, val int
	}

	columns := map[int][]entry{}
	queue := []entry{{root, 0, 0, root.Val}}
	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]
		columns[e.col] = append(columns[e.col], e)
		if e.node.Left != nil {
			queue = append(queue, entry{e.node.Left, e.col - 1, e.row + 1, e.node.Left.Val})
		}
		if e.node.Right != nil {
			queue = append(queue, entry{e.node.Right, e.col + 1, e.row + 1, e.node.Right.Val})
		}
	}

	columnKeys := make([]int, 0, len(columns))
	for key := range columns {
		columnKeys = append(columnKeys, key)
	}
	sort.Ints(columnKeys)

	var res [][]int
	for _, key := range columnKeys {
		entries := columns[key]
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].row != entries[j].row {
				return entries[i].row < entries[j].row
			}
			return entries[i].val < entries[j].val
		})
		vals := make([]int, len(entries))
		for i, e := range entries {
			vals[i] = e.val
		}
		res = append(res, vals)
	}
	return res
}
