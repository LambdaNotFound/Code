package tree

import . "gocode/types"

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
            res = append([]int{p.Val}, res...) // 1 Stack: [L, R, root]

            stack = append(stack, p)
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
