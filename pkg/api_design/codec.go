package apidesign

import (
	"fmt"
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
    // level ordering
    q := []*TreeNode{root}

    for len(q) > 0 {
        // pop out top
        top := q[0]
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
        // Left child
        if val, err := strconv.Atoi(slice[0]); err == nil {
            node.Left = &TreeNode{Val: val}
            q = append(q, node.Left)
        }
        slice = slice[1:]
        // Right child
        if len(slice) > 0 {
            if val, err := strconv.Atoi(slice[0]); err == nil {
                node.Right = &TreeNode{Val: val}
                q = append(q, node.Right)
            }
            slice = slice[1:]
        }
    }
    return head
}
