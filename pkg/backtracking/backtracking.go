package backtracking

import "sort"

/**
 * Backtracking (Recursive approach, DFS to try all possibilities)
 *
 *    for choice in choices {
 *        if notAllowed(choice) // Pruning => mark visited & skip if visited
 *            continue
 *
 *        makeChoice(choice)
 *        backtrack(path + choice, updatedChoices)
 *        undoChoice(choice) // backtrack
 *    }
 *
 * | Component              | Purpose                                       |
 * | ---------------------- | --------------------------------------------- |
 * | **Recursive function** | DFS-style traversal of the solution space     |
 * | **Base case**          | When a full solution is found, save it        |
 * | **Choices**            | The candidates to try at each step            |
 * | **Constraints check**  | Skip invalid paths early (pruning)            |
 * | **Backtrack (undo)**   | Remove the last choice before trying the next |
 *
 */

/**
 * 22. Generate Parentheses
 *
 * Given n pairs of parentheses, write a function to generate
 * all combinations of well-formed parentheses.
 *
 * backtrack(0, 0, "")
 * ├── backtrack(1, 0, "(")
 * │   ├── backtrack(2, 0, "((")
 * │   │   ├── backtrack(3, 0, "(((")
 * │   │   │   └── backtrack(3, 1, "((()")
 * │   │   │       └── backtrack(3, 2, "((())")
 * │   │   │           └── backtrack(3, 3, "((()))") ✓ append
 * │   │   └── backtrack(2, 1, "(()")
 * │   │       ├── backtrack(3, 1, "(()(")
 * │   │       │   └── backtrack(3, 2, "(()()")
 * │   │       │       └── backtrack(3, 3, "(()())") ✓ append
 * │   │       └── backtrack(2, 2, "(())")
 * │   │           └── backtrack(3, 2, "(())(")
 * │   │               └── backtrack(3, 3, "(())()") ✓ append
 * │   └── backtrack(1, 1, "()")
 * │       └── backtrack(2, 1, "()(")  ← right < left allows ) only after (
 * │           ├── backtrack(3, 1, "()((")
 * │           │   └── backtrack(3, 2, "()(()")
 * │           │       └── backtrack(3, 3, "()(())") ✓ append
 * │           └── backtrack(2, 2, "()()")
 * │               └── backtrack(3, 2, "()()(")
 * │                   └── backtrack(3, 3, "()()()") ✓ append
 *
 */
func generateParenthesis(n int) []string {
	result := make([]string, 0)
	path := make([]byte, 0, n*2) // pre-allocate exact capacity

	var dfs func(left, right int)
	dfs = func(left, right int) {
		// base case: all n pairs placed
		if left+right == n*2 {
			result = append(result, string(path))
			return
		}

		// place '(' if we still have opening parens left
		if left < n {
			path = append(path, '(')
			dfs(left+1, right)
			path = path[:len(path)-1]
		}

		// place ')' only if it won't exceed open count
		if right < left {
			path = append(path, ')')
			dfs(left, right+1)
			path = path[:len(path)-1]
		}
	}

	dfs(0, 0)
	return result
}

/**
 * 46. Permutations (ordered)
 *
 * Input: nums = [1,2,3]
 * Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
 *
 */
func permute(nums []int) [][]int {
	res := make([][]int, 0)
	visited := make([]bool, len(nums)) // index-based

	var backtrack func(path []int)
	backtrack = func(path []int) {
		if len(path) == len(nums) {
			res = append(res, append([]int{}, path...))
			return
		}

		for i := 0; i < len(nums); i++ {
			if visited[i] {
				continue
			}

			path = append(path, nums[i])
			visited[i] = true
			backtrack(path)
			visited[i] = false
			path = path[:len(path)-1]
		}
	}

	backtrack([]int{})
	return res
}

func permuteClaude(nums []int) [][]int {
	var res [][]int

	permutation := make([]int, len(nums))
	visited := make([]bool, len(nums))

	var backtrack func(int)
	backtrack = func(index int) {
		if index == len(nums) {
			copiedPermutation := make([]int, len(nums))
			copy(copiedPermutation, permutation)

			res = append(res, copiedPermutation)
			return
		}

		for i := 0; i < len(nums); i++ { // num[0] to num[i]
			if visited[i] == false {
				visited[i] = true
				permutation[index] = nums[i]
				backtrack(index + 1)
				visited[i] = false
			}
		}
	}

	backtrack(0)
	return res
}

func permuteWithSliceSpread(nums []int) [][]int {
	var result [][]int

	var backtrack func([]int, []int)
	backtrack = func(candidates, selected []int) {
		if len(candidates) == 0 {
			result = append(result, append([]int{}, selected...))
			return
		}
		for i, val := range candidates {
			sliceBefore := candidates[:i]
			newCandidates := append([]int{}, sliceBefore...)
			sliceAfter := candidates[i+1:]
			newCandidates = append(newCandidates, sliceAfter...)

			selected := append([]int{}, selected...)
			selected = append(selected, val)
			backtrack(newCandidates, selected)
		}
	}
	backtrack(nums, []int{})
	return result
}

/*
 *
 * 47. Permutations II
 *   !visited[i-1] means: "my identical left sibling is NOT currently in the path"
 *   → that sibling was already used at this position in a previous branch
 *   → skipping me prevents duplicate permutations
 *
 *   visited[i-1] means: "my identical left sibling IS currently in the path"
 *	 → sibling is at a DIFFERENT position than me
 *	 → picking me here creates a genuinely different permutation ✅
 */
func permuteUnique(nums []int) [][]int {
	sort.Ints(nums)
	res := make([][]int, 0)
	visited := make([]bool, len(nums)) // index-based

	var backtrack func(path []int)
	backtrack = func(path []int) {
		if len(path) == len(nums) {
			res = append(res, append([]int{}, path...))
			return
		}

		for i := 0; i < len(nums); i++ {
			if visited[i] {
				continue
			}

			if i > 0 && nums[i] == nums[i-1] && !visited[i-1] {
				continue
			}

			path = append(path, nums[i])
			visited[i] = true
			backtrack(path)
			visited[i] = false
			path = path[:len(path)-1]
		}
	}

	backtrack([]int{})
	return res
}

/*
 * 31. Next Permutation
 *
 * Narayana's algoritm
 */
// 1. Search max j: nums[j] > nums[j-1]
// 2. Search max l: l >= j && nums[l] > nums[j-1]
// 3. Swap nums[j-1], nums[l]
// 4. Reverse nums[j:]
// For example:
// 123654 -> 124653 -> 124356
func nextPermutation(nums []int) {
	j := len(nums) - 1
	for ; j > 0 && nums[j-1] >= nums[j]; j-- {
	}
	if j != 0 {
		l := len(nums) - 1
		for ; l > j-1 && nums[j-1] >= nums[l]; l-- {
		}
		nums[j-1], nums[l] = nums[l], nums[j-1]
	}
	for k := len(nums) - 1; j < k; j, k = j+1, k-1 {
		nums[j], nums[k] = nums[k], nums[j]
	}
}

/**
 * 77. Combinations (not ordered)
 */
func combine(n int, k int) [][]int {
	var res [][]int

	combination := make([]int, k)

	var backtrack func(int, int)
	backtrack = func(index int, num int) {
		if index == k {
			copiedCombination := make([]int, k)
			copy(copiedCombination, combination)

			res = append(res, copiedCombination)
			return
		}

		for i := num; i <= n; i++ {
			combination[index] = i
			backtrack(index+1, i+1)
		}
	}

	backtrack(0, 1)
	return res
}

func combineClaude(n int, k int) [][]int {
	res := make([][]int, 0)

	var backtrack func(start int, path []int)
	backtrack = func(start int, path []int) {
		if len(path) == k {
			res = append(res, append([]int{}, path...))
			return
		}

		// pruning: remaining elements must be enough to fill path to size k
		// need k-len(path) more elements, last valid start = n-(k-len(path))+1
		for i := start; i <= n-(k-len(path))+1; i++ {
			path = append(path, i)
			backtrack(i+1, path)
			path = path[:len(path)-1]
		}
	}

	backtrack(1, []int{})
	return res
}

/**
 * 39. Combination Sum
 *
 * Given an array of distinct integers candidates and a target integer target,
 * return a list of all unique combinations of candidates where the chosen numbers sum to target.
 *
 */
func combinationSum(candidates []int, target int) [][]int {
	var result [][]int

	var backtrack func(start, target int, selected []int)
	backtrack = func(start, target int, selected []int) {
		if target == 0 {
			result = append(result, append([]int{}, selected...))
			return
		}

		for i := start; i < len(candidates); i++ {
			val := candidates[i]
			if val > target {
				continue
			}
			// mutate and undo — no slice copies needed
			selected = append(selected, val)
			backtrack(i, target-val, selected)
			selected = selected[:len(selected)-1]
		}
	}

	backtrack(0, target, []int{})
	return result
}

/**
 * 377. Combination Sum IV
 *
 * Given an array of distinct integers nums and a target integer target, return the number of possible combinations that add up to target.
 *
 * dp[amount] += dp[amount - num]
 */
func combinationSum4(nums []int, target int) int {
	memo := make([]int, target+1)
	for i := range memo {
		memo[i] = -1
	}

	var dfs func(remaining int) int
	dfs = func(remaining int) int {
		if remaining < 0 {
			return 0
		}
		if remaining == 0 {
			return 1
		}
		if memo[remaining] != -1 {
			return memo[remaining]
		}

		count := 0
		for _, num := range nums {
			count += dfs(remaining - num)
		}
		memo[remaining] = count
		return count
	}

	return dfs(target)
}

/**
 * 78. Subsets
 *
 *  dfs(0, [])
 *  │   record []
 *  ├── i=0: path=[1]
 *  │   │   record [1]
 *  │   ├── i=1: path=[1,2]
 *  │   │   │   record [1,2]
 *  │   │   └── i=2: path=[1,2,3]
 *  │   │           record [1,2,3]
 *  │   │           return
 *  |   |
 *  │   └── i=2: path=[1,3]
 *  │           record [1,3]
 *  │           return
 *  |
 *  ├── i=1: path=[2]
 *  │   │   record [2]
 *  │   └── i=2: path=[2,3]
 *  │           record [2,3]
 *  │           return
 *  |
 *  └── i=2: path=[3]
 *          record [3]
 *          return
 */
func subsets(nums []int) [][]int {
	res := make([][]int, 0)

	var backtrack func(start int, path []int)
	backtrack = func(start int, path []int) {
		// record at every node — not just base case
		res = append(res, append([]int{}, path...))

		for i := start; i < len(nums); i++ {
			path = append(path, nums[i]) // make choice
			backtrack(i+1, path)
			path = path[:len(path)-1] // undo choice
		}
	}

	backtrack(0, []int{})
	return res
}

/**
 * 90. Subsets II
 */
func subsetsWithDup(nums []int) [][]int {
	sort.Ints(nums) // must sort first so duplicates are adjacent
	res := [][]int{}

	var backtrack func(start int, current []int)
	backtrack = func(start int, current []int) {
		res = append(res, append([]int{}, current...))

		for i := start; i < len(nums); i++ {
			// skip duplicate at the same recursion level
			if i > start && nums[i] == nums[i-1] {
				continue
			}
			current = append(current, nums[i])
			backtrack(i+1, current)
			current = current[:len(current)-1]
		}
	}

	backtrack(0, []int{})
	return res
}

/**
 * 17. Letter Combinations of a Phone Number
 */
func letterCombinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}

	phoneMap := map[byte]string{
		'2': "abc", '3': "def", '4': "ghi", '5': "jkl",
		'6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz",
	}

	res := []string{}
	var backtrack func(pos int, current []byte)
	backtrack = func(pos int, current []byte) {
		if pos == len(digits) {
			res = append(res, string(current))
			return
		}

		letters := phoneMap[digits[pos]]
		for i := range letters {
			current = append(current, letters[i])
			backtrack(pos+1, current)
			current = current[:len(current)-1]
		}
	}

	backtrack(0, []byte{})
	return res
}

func letterCombinationsBacktrack(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}

	dict := []string{"", "", "abc", "def", "ghi",
		"jkl", "mno", "pqrs", "tuv", "wxyz",
	}
	out := ""
	res := []string{} // pass in reference
	letterCombinationsHelper(digits, out, &res, dict)

	return res
}
func letterCombinationsHelper(digits string, out string, res *[]string, dict []string) {
	if len(digits) == 0 {
		*res = append(*res, out)
	} else {
		keys := dict[digits[0]-'0']
		next := digits[1:]
		for i := 0; i < len(keys); i++ {
			out = out + string(keys[i])
			letterCombinationsHelper(next, out, res, dict)
			out = out[:len(out)-1]
		}
	}
}
