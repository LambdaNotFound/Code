package backtracking

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
    res := []string{} // pass in reference
    letterCombinationsHelper(digits, out, &res, dict)

    return res
}

// in C++
// void letterCombinationsHelper(string digits, string& out, vector<string>& result, vector<string>& mapping) {
// golang pass by value
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
 */
func combinationSum(candidates []int, target int) [][]int {
    var result [][]int
    var backtrack func([]int, int, []int)

    backtrack = func(candidates []int, target int, selected []int) {
        if target == 0 {
            result = append(result, append([]int{}, selected...))
            return
        }
        for i, val := range candidates {
            if val <= target {
                newTarget := target - val
                selected := append([]int{}, selected...)
                selected = append(selected, val)
                newCandidates := append([]int{}, candidates[i:]...)
                backtrack(newCandidates, newTarget, selected)
            }
        }
    }
    backtrack(candidates, target, []int{})
    return result
}
