package interview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// buildTestTree constructs the 3-level tree used across evaluate tests:
//
//	root: X1 < 3
//	  left:  X2 < 1  → left="N", right="Y"
//	  right: X1 < 6  → left="N", right: X3 < 2 → left="Y", right="N"
func buildTestTree() *DecisionTree {
	dt := NewDecisionTree()
	root := dt.root

	left1, right1 := dt.AddSplit(root, "X1", 3)

	left2, right2 := dt.AddSplit(left1, "X2", 1)
	dt.SetLeafValue(left2, "N")
	dt.SetLeafValue(right2, "Y")

	left3, right3 := dt.AddSplit(right1, "X1", 6)
	dt.SetLeafValue(left3, "N")

	left4, right4 := dt.AddSplit(right3, "X3", 2)
	dt.SetLeafValue(left4, "Y")
	dt.SetLeafValue(right4, "N")

	return dt
}

func Test_DecisionTree_NewDecisionTree(t *testing.T) {
	dt := NewDecisionTree()
	assert.NotNil(t, dt)
	assert.NotNil(t, dt.root)
	assert.True(t, dt.root.isLeaf)
}

func Test_DecisionTree_AddSplit(t *testing.T) {
	dt := NewDecisionTree()
	root := dt.root

	left, right := dt.AddSplit(root, "X1", 5)

	assert.False(t, root.isLeaf)
	assert.Equal(t, "X1", root.signalName)
	assert.Equal(t, float64(5), root.constant)
	assert.NotNil(t, left)
	assert.NotNil(t, right)
	assert.True(t, left.isLeaf)
	assert.True(t, right.isLeaf)
}

func Test_DecisionTree_SetLeafValue(t *testing.T) {
	dt := NewDecisionTree()
	dt.SetLeafValue(dt.root, "Y")
	assert.Equal(t, "Y", dt.root.value)
}

func Test_DecisionTree_Evaluate(t *testing.T) {
	dt := buildTestTree()

	tests := []struct {
		name    string
		signals map[string]float64
		want    string
	}{
		{
			// X1=2 < 3 → left; X2=1 >= 1 → right → "Y"
			name:    "left branch right leaf",
			signals: map[string]float64{"X1": 2, "X2": 1, "X3": 11},
			want:    "Y",
		},
		{
			// X1=8 >= 3 → right; X1=8 >= 6 → right; X3=12 >= 2 → right → "N"
			name:    "right-right-right leaf",
			signals: map[string]float64{"X1": 8, "X2": 4, "X3": 12},
			want:    "N",
		},
		{
			// X1=1 < 3 → left; X2=0 < 1 → left → "N"
			name:    "left branch left leaf",
			signals: map[string]float64{"X1": 1, "X2": 0, "X3": 5},
			want:    "N",
		},
		{
			// X1=4 >= 3 → right; X1=4 < 6 → left → "N"
			name:    "right branch left leaf",
			signals: map[string]float64{"X1": 4, "X2": 0, "X3": 5},
			want:    "N",
		},
		{
			// X1=5 >= 3 → right; X1=5 < 6 → left → "N"
			name:    "right branch left leaf boundary",
			signals: map[string]float64{"X1": 5, "X2": 9, "X3": 1},
			want:    "N",
		},
		{
			// X1=7 >= 3 → right; X1=7 >= 6 → right; X3=1 < 2 → left → "Y"
			name:    "right-right-left leaf",
			signals: map[string]float64{"X1": 7, "X2": 9, "X3": 1},
			want:    "Y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, dt.Evaluate(tt.signals))
		})
	}
}

func Test_DecisionTree_Evaluate_SingleLeaf(t *testing.T) {
	dt := NewDecisionTree()
	dt.SetLeafValue(dt.root, "Y")
	assert.Equal(t, "Y", dt.Evaluate(map[string]float64{}))
}
