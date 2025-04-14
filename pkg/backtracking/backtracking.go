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
