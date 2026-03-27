package tree

import . "gocode/types"

/**
 * 105. Construct Binary Tree from Preorder and Inorder Traversal
 */
func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	// preorder[0] is always the root of the current subtree
	rootVal := preorder[0]

	// find root in inorder — splits into left and right subtrees
	idx := 0
	for ; idx < len(inorder); idx++ {
		if inorder[idx] == rootVal {
			break
		}
	}

	// idx = number of nodes in left subtree
	leftSize := idx
	root := &TreeNode{Val: rootVal}
	root.Left = buildTree(preorder[1:1+leftSize], inorder[:leftSize])
	root.Right = buildTree(preorder[1+leftSize:], inorder[leftSize+1:])

	return root
}

/**
 * 145. Binary Tree Postorder Traversal
 */
func postorderTraversal(root *TreeNode) []int {
	res, stack := []int{}, []*TreeNode{}
	p := root
	for p != nil || len(stack) != 0 {
		if p == nil {
			p = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			p = p.Left
		} else {
			res = append([]int{p.Val}, res...) // prepend to slice: [left, right, root
			stack = append(stack, p)           // push to stack: [root, right, left
			p = p.Right
		}
	}
	return res
}

/**
 * 144. Binary Tree Preorder Traversal
 */
func preorderTraversal(root *TreeNode) []int {
	res, stack := []int{}, []*TreeNode{}
	p := root
	for p != nil || len(stack) != 0 {
		if p == nil {
			p = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			p = p.Right // push right child if no left child
		} else {
			res = append(res, p.Val)

			stack = append(stack, p)
			p = p.Left // push left child to stack first
		}
	}
	return res
}

/**
 * 94. Binary Tree Inorder Traversal
 */
func inorderTraversal(root *TreeNode) []int {
	res, stack := []int{}, []*TreeNode{}
	p := root
	for p != nil || len(stack) != 0 {
		if p == nil {
			p = stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			res = append(res, p.Val)

			p = p.Right // push right child if no left child
		} else {
			stack = append(stack, p)
			p = p.Left // push left child to stack first
		}
	}
	return res
}

/**
 * 226. Invert Binary Tree
 *
 * 1. Recurisve approach: Time complexity: O(n), Space complexity: O(n)
 * 2. Iterative approach: Time complexity: O(n), Space complexity: O(n)
 */
func invertTree(root *TreeNode) *TreeNode {
	if root != nil {
		root.Left, root.Right = invertTree(root.Right), invertTree(root.Left)
	}
	return root
}

func invertTreeIterative(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	arr := []*TreeNode{root} // BFS

	for len(arr) != 0 {
		top := arr[len(arr)-1]
		arr = arr[:len(arr)-1]
		top.Left, top.Right = top.Right, top.Left
		if top.Left != nil {
			arr = append(arr, top.Left)
		}
		if top.Right != nil {
			arr = append(arr, top.Right)
		}
	}
	return root
}

/**
 * 662. Maximum Width of Binary Tree
 *
 */
func widthOfBinaryTree(root *TreeNode) int {
	if root == nil {
		return 0
	}

	type Node struct {
		node *TreeNode
		idx  int
	}

	queue := []Node{}
	queue = append(queue, Node{node: root, idx: 0})
	res := 0
	for len(queue) > 0 {
		size := len(queue)
		left, right := queue[0], queue[len(queue)-1]
		res = max(res, right.idx-left.idx+1)

		for i := 0; i < size; i++ {
			front := queue[0]
			queue = queue[1:]

			if front.node.Left != nil {
				queue = append(queue, Node{node: front.node.Left, idx: front.idx*2 + 1})
			}
			if front.node.Right != nil {
				queue = append(queue, Node{node: front.node.Right, idx: front.idx*2 + 2})
			}
		}
	}
	return res
}
