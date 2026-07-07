package math

/**
 * 169. Majority Element
 *
 * Given an array nums of size n, return the majority element.
 *
 * The majority element is the element that appears more than ⌊n / 2⌋ times.
 * You may assume that the majority element always exists in the array.
 *
 * Boyer-Moore Majority Vote Algorithm
 */
func majorityElement(nums []int) int {
    majority_element, majority_element_frequency := nums[0], 1
    for i := 1; i < len(nums); i++ {
        if majority_element_frequency == 0 {
            majority_element, majority_element_frequency = nums[i], 1
        } else {
            if nums[i] == majority_element {
                majority_element_frequency += 1
            } else {
                majority_element_frequency -= 1
            }
        }
    }
    return majority_element
}
