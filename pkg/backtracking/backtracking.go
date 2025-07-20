package backtracking

/**
 * Backtracking (Recursive approach, DFS to try all possibilities)
 *
 *    for choice in choices {
 *        if notAllowed(choice): continue
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
 * 46. Permutations
 */
func permute(nums []int) [][]int {
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

func permuteWithVisited(nums []int) [][]int {
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

        for i := 0; i < len(nums); i++ {
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

/**
 * 77. Combinations
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
            backtrack(index + 1, i + 1)
        }
    }

    backtrack(0, 1)
    return res
}

/**
 * 17. Letter Combinations of a Phone Number
 */
func letterCombinations(digits string) []string {
    if len(digits) == 0 {
        return []string{}
    }

    dict := []string{"", "", "abc", "def", "ghi",
        "jkl", "mno", "pqrs", "tuv", "wxyz",
    }
    out := ""
    res := []string{}

    var backtrack func(string, string)
    backtrack = func(str string, out string) {
        if len(str) == 0 {
            res = append(res, out)
        } else {
            keys := dict[str[0]-'0']
            next := str[1:]
            for i := 0; i < len(keys); i++ {
                out = out + string(keys[i])
                backtrack(next, out)
                out = out[:len(out)-1]
            }
        }
    }
    backtrack(digits, out)

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

/**
 * 39. Combination Sum
 *
 * Given an array of distinct integers candidates and a target integer target,
 * return a list of all unique combinations of candidates where the chosen numbers sum to target.
 *
 */
func combinationSum(candidates []int, target int) [][]int {
    var result [][]int
    var backtrack func(int, []int, []int)

    backtrack = func(target int, candidates, selected []int) {
        if target == 0 {
            result = append(result, append([]int{}, selected...))
            return
        }
        for i, val := range candidates {
            if val <= target {
                newTarget := target - val
                newCandidates := append([]int{}, candidates[i:]...)
                selected := append([]int{}, selected...)
                selected = append(selected, val)
                backtrack(newTarget, newCandidates, selected)
            }
        }
    }
    backtrack(target, candidates, []int{})
    return result
}

/**
 * 78. Subsets
 */
func subsets(nums []int) [][]int {
    var result [][]int
    var curr []int

    var search func(int)
    search = func(index int) {
        if index == len(nums) {
            subset := make([]int, len(curr))
            copy(subset, curr)
            result = append(result, subset)
            return
        }

        curr = append(curr, nums[index])
        search(index + 1)

        curr = curr[:len(curr)-1]
        search(index + 1)
    }

    search(0)
    return result
}
