package solid_coding

/*
 * 339. Nested List Weight Sum
 *
 * Example 1:
 * Input:  [[1,1],2,[1,1]]
 * Output: 10
 * Explanation: Four 1's at depth 2, one 2 at depth 1.
 *              1*2 + 1*2 + 2*1 + 1*2 + 1*2 = 10
 *
 */

type NestedInteger struct{}

func (n NestedInteger) IsInteger() bool
func (n NestedInteger) GetInteger() int
func (n NestedInteger) GetList() []*NestedInteger

func depthSum(nestedList []*NestedInteger) int {
	queue := make([]*NestedInteger, len(nestedList))
	copy(queue, nestedList)

	res, depth := 0, 1
	for len(queue) > 0 {
		size := len(queue) // snapshot current level
		for i := 0; i < size; i++ {
			cur := queue[0]
			queue = queue[1:]

			if cur.IsInteger() {
				res += cur.GetInteger() * depth // multiply by current depth
			} else {
				queue = append(queue, cur.GetList()...)
			}
		}
		depth++
	}

	return res
}

func depthSumDFS(nestedList []*NestedInteger) int {
	var dfs func(list []*NestedInteger, depth int) int
	dfs = func(list []*NestedInteger, depth int) int {
		total := 0
		for _, item := range list {
			if item.IsInteger() {
				total += item.GetInteger() * depth
			} else {
				total += dfs(item.GetList(), depth+1)
			}
		}
		return total
	}

	return dfs(nestedList, 1)
}
