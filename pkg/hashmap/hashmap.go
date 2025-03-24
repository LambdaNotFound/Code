package hashmap

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
