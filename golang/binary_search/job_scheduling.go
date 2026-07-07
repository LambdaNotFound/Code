package binarysearch

import "sort"

/**
 * 1235. Maximum Profit in Job Scheduling
 *
 * dp[i] = max(dp[i-1], jobs[i].profit + dp[findLastNonOverlapping(i)])
 *
 * top-down DFS needs start-time order, bottom-up DP needs end-time order.
 */
func jobScheduling(startTime []int, endTime []int, profit []int) int {
	n := len(startTime)
	type Job struct {
		start, end, profit int
	}
	jobs := make([]Job, n)
	for i := range startTime {
		jobs[i] = Job{start: startTime[i], end: endTime[i], profit: profit[i]}
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].end < jobs[j].end
	})

	findPrevJob := func(pos int) int {
		left, right := 0, pos
		for left < right {
			mid := left + (right-left)/2
			if jobs[mid].end <= jobs[pos].start {
				left = mid + 1
			} else {
				right = mid
			}
		}
		return left - 1
	}

	dp := make([]int, n)
	dp[0] = jobs[0].profit
	for i := 1; i < n; i++ {
		prev := findPrevJob(i)
		prevProfit := 0
		if prev >= 0 {
			prevProfit = dp[prev]
		}
		dp[i] = max(dp[i-1], prevProfit+jobs[i].profit)
	}

	return dp[n-1]
}

func jobSchedulingTopDown(startTime []int, endTime []int, profit []int) int {
	type job struct {
		start, end, profit int
	}
	jobs := make([]job, len(startTime))
	for i := range startTime {
		jobs[i] = job{startTime[i], endTime[i], profit[i]}
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].start < jobs[j].start
	})

	// find first job index where start >= endTime
	findNextJob := func(pos int) int {
		left, right := pos+1, len(jobs) // ← start from pos+1
		for left < right {
			mid := left + (right-left)/2
			if jobs[pos].end <= jobs[mid].start {
				right = mid
			} else {
				left = mid + 1
			}
		}
		return left
	}

	memo := make(map[int]int)
	var dfs func(int) int
	dfs = func(pos int) int {
		if pos == len(jobs) {
			return 0
		}
		if val, found := memo[pos]; found {
			return val
		}

		// include current job: add its profit + best from next valid job
		include := jobs[pos].profit + dfs(findNextJob(pos))
		// exclude current job: best from next job
		exclude := dfs(pos + 1)

		memo[pos] = max(include, exclude)
		return memo[pos]
	}

	return dfs(0)
}
