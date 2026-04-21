package apidesign

import "sort"

/**
 * 981. Time Based Key-Value Store
 *
 * Implement the TimeMap class:
 *
 * TimeMap() Initializes the object of the data structure.
 * void set(String key, String value, int timestamp) Stores the key
 * key with the value value at the given time timestamp.
 *
 * String get(String key, int timestamp) Returns a value such that
 * set was called previously, with timestamp_prev <= timestamp.
 * If there are multiple such values, it returns the value associated with
 * the largest timestamp_prev. If there are no values, it returns "".
 *
 */
type TimeMap struct {
	store map[string][]entry // Map → store: lowercase, descriptive
}

type entry struct { // Pair → entry: domain-specific name
	timestamp int    // Time → timestamp: matches the problem's language
	value     string // Val → value: full word
}

func ConstructorTimeMap() TimeMap {
	return TimeMap{
		store: make(map[string][]entry),
	}
}

func (tm *TimeMap) Set(key string, value string, timestamp int) {
	tm.store[key] = append(tm.store[key], entry{timestamp: timestamp, value: value})
}

func (tm *TimeMap) Get(key string, timestamp int) string {
	entries, found := tm.store[key]
	if !found {
		return ""
	}

	left, right := 0, len(entries)
	for left < right {
		mid := left + (right-left)/2
		if timestamp < entries[mid].timestamp {
			right = mid
		} else {
			left = mid + 1
		}
	}

	// left is now the insertion point — the answer is at left-1
	if left == 0 {
		return "" // all timestamps are greater than query
	}
	return entries[left-1].value
}

/*
 * // lower bound (original sort.Search)
 * sort.Search(len(arr), func(i int) bool {
 *    return arr[i].timestamp >= timestamp  // lands ON target
 * })
 *
 * // upper bound (your version)
 * sort.Search(len(arr), func(i int) bool {
 *    return arr[i].timestamp > timestamp   // lands PAST target
 * })
 */

func (tm *TimeMap) GetByUpperBound(key string, timestamp int) string {
	if _, found := tm.store[key]; !found {
		return ""
	}

	arr := tm.store[key]
	index := sort.Search(len(arr), func(i int) bool {
		return arr[i].timestamp > timestamp
	})

	if index == 0 {
		return "" // target < all timestamps
	}
	return arr[index-1].value // index-1 is the floor
}
