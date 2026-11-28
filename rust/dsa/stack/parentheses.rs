pub struct Solution;

impl Solution {
    /**
     * 32. Longest Valid Parentheses
     *  1. stack
     *  2. DP
     *  3. two pointers
     *
     * translated from go/dsa/stack/parentheses.go
     */
    pub fn longest_valid_parentheses(s: String) -> i32 {
        let mut stack: Vec<i32> = vec![-1]; // boundary: track the invalid indices
        let mut res = 0;

        for (i, c) in s.bytes().enumerate() {
            let i = i as i32;
            if c == b'(' {
                stack.push(i); // push index of '('
            } else {
                stack.pop(); // pop the matching '(' or the boundary

                match stack.last() {
                    None => stack.push(i),                // new boundary: unmatched ')'
                    Some(&top) => res = res.max(i - top), // valid length
                }
            }
        }
        res
    }

    /**
     * Case 1: s[i-1] == '('   → direct match "()"
     *         dp[i] = dp[i-2] + 2
     *
     * Case 2: s[i-1] == ')'   → previous char closes its own group "))"
     *     need to check if char BEFORE that group matches
     *     j = i - dp[i-1] - 1
     *     if s[j] == '(' :
     *         dp[i] = dp[i-1] + 2 + dp[j-1]
     */
    pub fn longest_valid_parentheses_dp(s: String) -> i32 {
        let s = s.as_bytes();
        let mut dp = vec![0usize; s.len()];
        let mut res = 0;

        for i in 1..s.len() {
            if s[i] == b')' {
                if s[i - 1] == b'(' {
                    // case 1: direct "()" match
                    dp[i] = if i >= 2 { dp[i - 2] } else { 0 } + 2;
                } else if s[i - 1] == b')' {
                    // case 2: ends with "))" — look before inner group;
                    // checked_sub stands in for Go's j >= 0
                    if let Some(j) = i.checked_sub(dp[i - 1] + 1)
                        && s[j] == b'('
                    {
                        dp[i] = dp[i - 1] + 2 + if j > 0 { dp[j - 1] } else { 0 };
                    }
                }
                res = res.max(dp[i]);
            }
        }
        res as i32
    }

    /**
     * Count '(' and ')'; a run is valid when the counts are equal. Left to right
     * resets on an unmatched ')', but never sees the counts equal inside "(()",
     * so a second pass right to left, resetting on an unmatched '(', finds it.
     */
    pub fn longest_valid_parentheses_two_pointers(s: String) -> i32 {
        let mut res = 0;

        // pass 1: left → right
        let (mut opens, mut closes) = (0, 0);
        for c in s.bytes() {
            if c == b'(' {
                opens += 1;
            } else {
                closes += 1;
            }
            if opens == closes {
                res = res.max(2 * closes);
            } else if closes > opens {
                (opens, closes) = (0, 0); // unmatched ')' — reset
            }
        }

        // pass 2: right → left
        (opens, closes) = (0, 0);
        for c in s.bytes().rev() {
            if c == b'(' {
                opens += 1;
            } else {
                closes += 1;
            }
            if opens == closes {
                res = res.max(2 * opens);
            } else if opens > closes {
                (opens, closes) = (0, 0); // unmatched '(' — reset
            }
        }

        res
    }

    /**
     * 20. Valid Parentheses
     *
     * push the closer each opener expects; a closer must match the top
     * translated from go/dsa/stack/parentheses.go
     */
    pub fn is_valid(s: String) -> bool {
        let mut stack: Vec<u8> = Vec::with_capacity(s.len());
        for c in s.bytes() {
            match c {
                b'(' => stack.push(b')'),
                b'[' => stack.push(b']'),
                b'{' => stack.push(b'}'),
                // closing bracket
                _ => {
                    if stack.pop() != Some(c) {
                        return false;
                    }
                }
            }
        }
        stack.is_empty()
    }
}

#[cfg(test)]
mod tests {
    use super::Solution;

    #[test]
    fn longest_valid_parentheses() {
        // same cases as Test_longestValidParentheses in
        // go/dsa/stack/stack_extra_test.go, then extra edge cases
        let cases: &[(&str, &str, i32)] = &[
            ("leetcode_1", "(()", 2),
            ("leetcode_2", ")()())", 4),
            ("empty", "", 0),
            ("all_open", "(((", 0),
            ("all_close", ")))", 0),
            ("full_match", "()()", 4),
            ("nested", "(())", 4),
            ("complex", "()(())", 6),
            ("trailing_open", "()((", 2),
            ("single_pair", "()", 2),
            ("interleaved", "())()(", 2),
            ("middle_run", "(()))())(", 4),
            ("run_after_reset", ")()())()()(", 4),
            ("split_by_open", "()(()", 2),
            ("unmatched_open_left", "(()(((()", 2),
            ("joined_groups", ")(((((()())()()))()(()))(", 22),
        ];
        type Method = fn(String) -> i32;
        let methods: &[(&str, Method)] = &[
            ("stack", Solution::longest_valid_parentheses),
            ("dp", Solution::longest_valid_parentheses_dp),
            (
                "two_pointers",
                Solution::longest_valid_parentheses_two_pointers,
            ),
        ];
        for &(name, s, expected) in cases {
            for &(method, f) in methods {
                assert_eq!(f(s.to_string()), expected, "case {name}, {method}");
            }
        }
    }

    #[test]
    fn is_valid() {
        // same cases as Test_isValid in go/dsa/stack/stack_test.go, plus three
        let cases: &[(&str, &str, bool)] = &[
            ("empty", "", true),
            ("single_pair", "()", true),
            ("multiple_types", "()[]{}", true),
            ("wrong_order", "(]", false),
            ("nested", "([{}])", true),
            ("only_open", "(((", false),
            ("only_close", ")))", false),
            ("mismatch", "([)]", false),
            ("close_before_open", "){", false),
            ("extra_close", "())", false),
            ("extra_open", "(()", false),
        ];
        for &(name, s, expected) in cases {
            assert_eq!(Solution::is_valid(s.to_string()), expected, "case {name}");
        }
    }
}
