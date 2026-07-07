package interview


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
