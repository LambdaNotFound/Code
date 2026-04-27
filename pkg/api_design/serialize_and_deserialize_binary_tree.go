package apidesign

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	. "gocode/types"
)

/**
 * 297. Serialize and Deserialize Binary Tree
 */
type Codec struct {
}

func Constructor() Codec {
	return Codec{}
}

// Serializes a tree to a single string.
func (this *Codec) serialize(root *TreeNode) string {
	if root == nil {
		return ""
	}
	res := []string{}
	q := []*TreeNode{root} // BFS

	for len(q) > 0 {
		top := q[0] // append str on dequeue
		q = q[1:]
		if top == nil {
			res = append(res, "nil")
		} else {
			res = append(res, fmt.Sprintf("%d", top.Val))
			q = append(q, top.Left)
			q = append(q, top.Right)
		}
	}
	return strings.Join(res, ",")
}

// Deserializes your encoded data to tree.
func (this *Codec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}
	slice := strings.Split(data, ",")
	val, _ := strconv.Atoi(slice[0])
	slice = slice[1:]

	head := &TreeNode{Val: val}
	q := []*TreeNode{head}
	for len(slice) > 0 {
		node := q[0]
		q = q[1:]

		if len(slice) > 0 { // Left child, err != nil when parsing "nil"
			if val, err := strconv.Atoi(slice[0]); err == nil {
				node.Left = &TreeNode{Val: val}
				q = append(q, node.Left)
			}
			slice = slice[1:]
		}

		if len(slice) > 0 { // Right child, err != nil when parsing "nil"
			if val, err := strconv.Atoi(slice[0]); err == nil {
				node.Right = &TreeNode{Val: val}
				q = append(q, node.Right)
			}
			slice = slice[1:]
		}
	}
	return head
}

func (this *Codec) serializeClaude(root *TreeNode) string {
	tokens := []string{}
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node != nil {
			tokens = append(tokens, strconv.Itoa(node.Val))
			queue = append(queue, node.Left)
			queue = append(queue, node.Right)
		} else {
			tokens = append(tokens, "null")
		}
	}
	return strings.Join(tokens, ",")
}

func (this *Codec) deserializeClaude(data string) *TreeNode {
	tokens := strings.Split(data, ",")
	idx := 0

	readNext := func() *TreeNode {
		if idx >= len(tokens) || tokens[idx] == "null" {
			idx++
			return nil
		}
		val, _ := strconv.Atoi(tokens[idx])
		idx++
		return &TreeNode{Val: val}
	}

	root := readNext()
	if root == nil {
		return nil
	}

	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		node.Left = readNext()
		node.Right = readNext()

		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	return root
}

/**
 * 105. Construct Binary Tree from Preorder and Inorder Traversal
 *
 * 1. Recursive Approach
 *
 */
func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	i := slices.Index(inorder, preorder[0])
	lpre, rpre := preorder[1:i+1], preorder[i+1:]
	lin, rin := inorder[:i], inorder[i+1:]

	return &TreeNode{
		Val:   preorder[0],
		Left:  buildTree(lpre, lin),
		Right: buildTree(rpre, rin),
	}
}
