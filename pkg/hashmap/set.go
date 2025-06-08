package hashmap

/**
 * 217. Contains Duplicate
 */
func containsDuplicate(nums []int) bool {
    set := make(map[int]bool) // map[int]struct{}
    for _, num := range nums {
        if _, hasNum := set[num]; hasNum {
            return true
        }
        set[num] = true // struct{}{}
    }
    return false
}
