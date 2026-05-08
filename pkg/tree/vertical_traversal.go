package tree

import (
	. "gocode/types"
	"sort"
)

/**
 * 987. Vertical Order Traversal of a Binary Tree
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
		val, row, col int
	}

	var items []entry
	queue := []entry{{root, root.Val, 0, 0}}
	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]
		items = append(items, e)
		if e.node.Left != nil {
			queue = append(queue, entry{e.node.Left, e.node.Left.Val, e.row + 1, e.col - 1})
		}
		if e.node.Right != nil {
			queue = append(queue, entry{e.node.Right, e.node.Right.Val, e.row + 1, e.col + 1})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].col != items[j].col {
			return items[i].col < items[j].col
		}
		if items[i].row != items[j].row {
			return items[i].row < items[j].row
		}
		return items[i].val < items[j].val
	})

	var res [][]int
	for i := 0; i < len(items); {
		col := items[i].col
		var group []int
		for i < len(items) && items[i].col == col {
			group = append(group, items[i].val)
			i++
		}
		res = append(res, group)
	}
	return res
}
