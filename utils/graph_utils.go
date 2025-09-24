package utils

import . "gocode/types"

// GraphsEqual checks structural equality of two graphs
func GraphsEqual(a, b *Node, visited map[*Node]*Node) bool {
    if a == nil || b == nil {
        return a == b
    }
    if a.Val != b.Val {
        return false
    }
    if v, ok := visited[a]; ok {
        return v == b
    }
    visited[a] = b
    if len(a.Neighbors) != len(b.Neighbors) {
        return false
    }
    for i := range a.Neighbors {
        if !GraphsEqual(a.Neighbors[i], b.Neighbors[i], visited) {
            return false
        }
    }
    return true
}