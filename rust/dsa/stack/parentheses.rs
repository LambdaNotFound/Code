pub struct Solution;

impl Solution {
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
