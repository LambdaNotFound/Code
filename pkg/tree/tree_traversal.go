package tree

import . "gocode/types"

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

            p = p.Right // right child at last
        } else {
            stack = append(stack, p)
            p = p.Left
        }
    }

    return res
}
