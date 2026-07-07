package array

import "sort"

/**
 * 41. First Missing Positive
 *
 * in O(n) time and uses O(1) auxiliary space.
 *
 * Input: nums = [1,2,0]
 * Output: 3
 *
 * Input: nums = [3,4,-1,1]
 * Output: 2
 *
 * goal: nums[i] == i+1
 * move all nums[i] to nums[nums[i]-1]
 *
 */
func firstMissingPositive(nums []int) int {
	n := len(nums)
	for i := range nums { // nums[i] != nums[nums[i]-1], Input: [1,1] => infinite loop
		for nums[i]-1 >= 0 && nums[i]-1 < n && nums[i] != nums[nums[i]-1] {
			nums[i], nums[nums[i]-1] = nums[nums[i]-1], nums[i]
		}
	}

	for i := range nums {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return n + 1
}

/**
 * 80. Remove Duplicates from Sorted Array II
 *
 * each unique element appears at most twice
 * Input: nums = [1,1,1,2,2,3]
 * Output: 5, nums = [1,1,2,2,3,_]
 */
func removeDuplicatesFromSortedArray(nums []int) int {
	index, count := 1, 1
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[index-1] {
			count++
		} else {
			count = 1
		}

		if count <= 2 {
			nums[index] = nums[i]
			index++
		}
	}
	return index
}

/**
 * Team Arrangement
 *
 * You are given two arrays of integers, heights1 and heights2, representing the heights of players from two different teams. Your task is to determine if it's possible to arrange the teams for a photograph according to the following rules:
 * One team must stand entirely in front of the other.
 * Every player in the front row must be strictly shorter than the player standing directly behind them.
 * Players can be rearranged within their own team's row.
 * The team with fewer players must stand in the front row. If the teams are of equal size, either team can be the front row.
 * Your function should return True if such an arrangement is possible, and False otherwise.
 *
 * BBBBBBBB
 * FF FFF     ok
 *   FFFF FF  !ok
 */
func canArrange(team1, team2 []int) bool {
	sort.Ints(team1)
	sort.Ints(team2)

	check := func(front, back []int) bool {
		fi := 0
		for bi := 0; fi < len(front) && bi < len(back); bi++ {
			if front[fi] < back[bi] {
				fi++
			}
		}
		return fi == len(front)
	}

	if len(team1) > len(team2) {
		team1, team2 = team2, team1
	}
	if len(team1) < len(team2) {
		return check(team1, team2)
	}
	return check(team1, team2) || check(team2, team1)
}
