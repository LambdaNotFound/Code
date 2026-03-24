package apidesign

import "math/rand"

/**
 * 380. Insert Delete GetRandom O(1)
 *
 * hashmap + array
 */

type RandomizedSet struct {
	nums     []int
	indexMap map[int]int // val → index in nums
}

func ConstructorRandomizedSet() RandomizedSet {
	return RandomizedSet{
		nums:     make([]int, 0),
		indexMap: make(map[int]int),
	}
}

func (rs *RandomizedSet) Insert(val int) bool {
	if _, exists := rs.indexMap[val]; exists {
		return false
	}
	rs.nums = append(rs.nums, val)
	rs.indexMap[val] = len(rs.nums) - 1
	return true
}

func (rs *RandomizedSet) Remove(val int) bool {
	idx, exists := rs.indexMap[val]
	if !exists {
		return false
	}

	// swap val with last element
	last := rs.nums[len(rs.nums)-1]
	rs.nums[idx] = last
	rs.indexMap[last] = idx

	// remove last element
	rs.nums = rs.nums[:len(rs.nums)-1]
	delete(rs.indexMap, val)

	return true
}

func (rs *RandomizedSet) GetRandom() int {
	return rs.nums[rand.Intn(len(rs.nums))]
}
