package twopointers

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
