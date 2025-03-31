package hashmap

/**
 * 1. Two Sum
 */
func twoSum(nums []int, target int) []int {
    hashmap := make(map[int]int)
    // result := make([]int, 0)

    for i, val := range nums {
        diff := target - val
        if index, exist := hashmap[diff]; exist {
            // result = append(result, index, i)
            return []int{index, i}
        } else {
            hashmap[val] = i
        }
    }

    return nil
}

/**
 * 383. Ransom Note
 */
func canConstruct(ransomNote string, magazine string) bool {
    hashmap := make(map[rune]int)
    for _, val := range magazine {
        hashmap[val] += 1
    }

    for _, val := range ransomNote {
        if hashmap[val] > 0 {
            hashmap[val] -= 1
        } else {
            return false
        }
    }
    return true
}
