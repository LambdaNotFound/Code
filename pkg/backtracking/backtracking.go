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
    res := []string{}

    var backtrack func(string, string)
    backtrack = func(target string, out string) {
        if len(target) == 0 {
            res = append(res, out)
        } else {
            keys := dict[target[0]-'0']
            next := target[1:]
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

    backtrack = func(target int, candidates []int, selected []int) {
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
