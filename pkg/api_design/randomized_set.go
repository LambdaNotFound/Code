package apidesign

import "math/rand"

/**
 * 380. Insert Delete GetRandom O(1)
 *
 * hashmap + array
 */

type RandomizedSet struct {
	nums    []int
	indices map[int]int // val → index in nums
}

func ConstructorRandomizedSet() RandomizedSet {
	return RandomizedSet{
		nums:    make([]int, 0),
		indices: make(map[int]int),
	}
}

func (rs *RandomizedSet) Insert(val int) bool {
	if _, exists := rs.indices[val]; exists {
		return false
	}
	rs.nums = append(rs.nums, val)
	rs.indices[val] = len(rs.nums) - 1
	return true
}

func (rs *RandomizedSet) Remove(val int) bool {
	idx, exists := rs.indices[val]
	if !exists {
		return false
	}

	// swap val with last element
	last := rs.nums[len(rs.nums)-1]
	rs.nums[idx] = last
	rs.indices[last] = idx

	// remove last element
	rs.nums = rs.nums[:len(rs.nums)-1]
	delete(rs.indices, val)

	return true
}

func (rs *RandomizedSet) GetRandom() int {
	return rs.nums[rand.Intn(len(rs.nums))]
}

/**
 * 381. Insert Delete GetRandom O(1) - Duplicates allowed
 *
 * hashmap -> hashset + array
 */
type RandomizedCollection struct {
	array   []int
	indices map[int]map[int]bool
}

func ConstructorRandomizedCollection() RandomizedCollection {
	return RandomizedCollection{
		indices: make(map[int]map[int]bool),
	}
}

func (r *RandomizedCollection) Insert(val int) bool {
	idx := len(r.array)
	r.array = append(r.array, val)

	if r.indices[val] == nil {
		r.indices[val] = make(map[int]bool)
	}
	r.indices[val][idx] = true

	return len(r.indices[val]) == 1
}

func (r *RandomizedCollection) Remove(val int) bool {
	if len(r.indices[val]) == 0 {
		return false
	}

	lastPos := len(r.array) - 1
	lastVal := r.array[lastPos]

	// Pick any index where val lives
	var idx int
	for idx = range r.indices[val] {
		break
	}

	// If lastVal == val, always pick lastPos to avoid orphaning other indices
	if lastVal == val {
		idx = lastPos
	}

	// Only swap if val is not already the last element
	if idx != lastPos {
		r.array[idx] = lastVal
		r.indices[lastVal][idx] = true
		delete(r.indices[lastVal], lastPos)
	}

	// Remove val's index entry
	delete(r.indices[val], idx)
	if len(r.indices[val]) == 0 {
		delete(r.indices, val)
	}

	r.array = r.array[:lastPos]
	return true
}

func (r *RandomizedCollection) GetRandom() int {
	return r.array[rand.Intn(len(r.array))]
}
