package interview

import "fmt"

type Node struct {
	value      string
	signalName string
	constant   float64

	left  *Node
	right *Node

	isLeaf bool
}

type DecisionTree struct {
	root *Node
}

func NewDecisionTree() *DecisionTree {
	// Start as a single leaf with no value set yet
	return &DecisionTree{root: &Node{isLeaf: true}}
}

// AddSplit converts a leaf into an interior node and returns the two new leaves.
// The caller must hold onto the returned leaves to split them further later.
func (dt *DecisionTree) AddSplit(leaf *Node, signalName string, constant float64) (*Node, *Node) {
	// Convert this leaf into an interior node in-place.
	// We mutate the node directly so any existing pointer to it stays valid.
	leaf.isLeaf = false
	leaf.signalName = signalName
	leaf.constant = constant

	leaf.left = &Node{isLeaf: true}
	leaf.right = &Node{isLeaf: true}

	return leaf.left, leaf.right
}

// SetLeafValue assigns a return value to a leaf node.
func (dt *DecisionTree) SetLeafValue(leaf *Node, value string) {
	leaf.value = value
}

// Evaluate traverses the tree using the provided signals and returns the leaf value.
func (dt *DecisionTree) Evaluate(signals map[string]float64) string {
	curr := dt.root
	for !curr.isLeaf {
		// Route left if signal < constant, right if signal >= constant
		if signals[curr.signalName] < curr.constant {
			curr = curr.left
		} else {
			curr = curr.right
		}
	}
	return curr.value
}

func testDecisionTree() {
	dt := NewDecisionTree()
	root := dt.root

	// Split root on X1 < 3
	left1, right1 := dt.AddSplit(root, "X1", 3)

	// Left branch: split on X2 < 1
	left2, right2 := dt.AddSplit(left1, "X2", 1)
	dt.SetLeafValue(left2, "N")
	dt.SetLeafValue(right2, "Y")

	// Right branch: split on X1 < 6
	left3, right3 := dt.AddSplit(right1, "X1", 6)
	dt.SetLeafValue(left3, "N")

	// Right-right branch: split on X3 < 2
	left4, right4 := dt.AddSplit(right3, "X3", 2)
	dt.SetLeafValue(left4, "Y")
	dt.SetLeafValue(right4, "N")

	// --- Test cases ---

	// X1=2 < 3 → left; X2=1 >= 1 → right → "Y"
	fmt.Println(dt.Evaluate(map[string]float64{"X1": 2, "X2": 1, "X3": 11})) // Y

	// X1=8 >= 3 → right; X1=8 >= 6 → right; X3=12 >= 2 → right → "N"
	fmt.Println(dt.Evaluate(map[string]float64{"X1": 8, "X2": 4, "X3": 12})) // N

	// X1=1 < 3 → left; X2=0 < 1 → left → "N"
	fmt.Println(dt.Evaluate(map[string]float64{"X1": 1, "X2": 0, "X3": 5})) // N

	// X1=4 >= 3 → right; X1=4 < 6 → left → "N"
	fmt.Println(dt.Evaluate(map[string]float64{"X1": 4, "X2": 0, "X3": 5})) // N

	// X1=5 >= 3 → right; X1=5 < 6 → left → "N"
	fmt.Println(dt.Evaluate(map[string]float64{"X1": 5, "X2": 9, "X3": 1})) // N

	// X1=7 >= 3 → right; X1=7 >= 6 → right; X3=1 < 2 → left → "Y"
	fmt.Println(dt.Evaluate(map[string]float64{"X1": 7, "X2": 9, "X3": 1})) // Y
}
