package recursion

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

/*
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

// 341. Flatten Nested List Iterator

type NestedIterator struct {
	items   []int
	current int
}

func Constructor(nestedList []*NestedInteger) *NestedIterator {
	items := []int{}
	var dfs func(list []*NestedInteger)
	dfs = func(list []*NestedInteger) {
		for _, item := range list {
			if item.IsInteger() {
				items = append(items, item.GetInteger())
			} else {
				dfs(item.GetList())
			}
		}
	}
	dfs(nestedList)

	return &NestedIterator{
		items: items,
	}
}

func (it *NestedIterator) Next() int {
	val := it.items[it.current]
	it.current++
	return val
}

func (it *NestedIterator) HasNext() bool {
	return it.current < len(it.items)
}
*/
