use std::collections::HashMap;

pub struct Solution;

impl Solution {
    /**
     * 1. Two Sum
     *
     * one-pass, to avoid using same number twice
     * translated from go/dsa/hash_table/hashmap.go; Go's nil becomes an empty Vec
     */
    pub fn two_sum(nums: Vec<i32>, target: i32) -> Vec<i32> {
        let mut complements: HashMap<i32, i32> = HashMap::with_capacity(nums.len());
        for (i, &num) in nums.iter().enumerate() {
            if let Some(&index) = complements.get(&(target - num)) {
                return vec![index, i as i32];
            }
            complements.insert(num, i as i32);
        }
        vec![]
    }
}

#[cfg(test)]
mod tests {
    use super::Solution;

    #[test]
    fn two_sum() {
        // same cases as Test_twoSum in go/dsa/hash_table/hashmap_test.go
        let cases: &[(&str, &[i32], i32, &[i32])] = &[
            ("basic_case", &[2, 7, 11, 15], 9, &[0, 1]),
            ("another_case", &[3, 2, 4], 6, &[1, 2]),
            ("repeated_numbers", &[3, 3], 6, &[0, 1]),
            ("larger_array", &[1, 2, 3, 4, 5, 6], 11, &[4, 5]),
            ("negatives", &[-1, -2, -3, -4, -5], -8, &[2, 4]),
            ("target_zero", &[-3, 3, 1], 0, &[0, 1]),
            ("answer_at_end", &[1, 2, 3, 4], 7, &[2, 3]),
            ("zero_and_value", &[0, 4, 3, 0], 0, &[0, 3]),
            ("no_solution", &[1, 2, 3], 10, &[]),
        ];
        for &(name, nums, target, expected) in cases {
            assert_eq!(
                Solution::two_sum(nums.to_vec(), target),
                expected,
                "case {name}"
            );
        }
    }
}
