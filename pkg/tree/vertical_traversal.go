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

	type entry struct {
		node     *TreeNode
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
