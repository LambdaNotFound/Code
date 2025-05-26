package dynamic_programming

import "sort"

/**
 * 1235. Maximum Profit in Job Scheduling
 *
 * We have n jobs, where every job is scheduled to be done from startTime[i] to endTime[i],
 * obtaining a profit of profit[i].
 *
 * You're given the startTime, endTime and profit arrays, return the maximum profit you can take
 * such that there are no two jobs in the subset with overlapping time range.
 *
 * 1. DFS w/ memo
 *
 */
func jobScheduling_dfs(startTime []int, endTime []int, profit []int) int {
    // initialize jobs array & sort it by start time
    type Job struct {
        start, end, profit int
    }
    jobs := make([]Job, len(startTime))
    for i := 0; i < len(startTime); i++ {
        jobs[i] = Job{startTime[i], endTime[i], profit[i]}
    }
    sort.Slice(jobs, func(i, j int) bool { return jobs[i].start < jobs[j].start })

    // binary search to find the next non-conflicting job
    searchNextJob := func(i int) int {
        left, right := i+1, len(jobs)-1
        for left <= right {
            mid := (left + right) / 2
            if jobs[mid].start >= jobs[i].end {
                left, right = left, mid-1
            } else {
                left, right = mid+1, right
            }
        }
        return left
    }

    memo := make(map[int]int)
    var dfs func(int) int
    dfs = func(pos int) int {
        if pos == len(jobs) { // base condition
            return 0
        }

        if ans, ok := memo[pos]; ok { // memoized branch
            return ans
        }

        res := jobs[pos].profit // take this job to execute
        res += dfs(searchNextJob(pos))
        alt := dfs(pos + 1) // not take this job
        res = max(res, alt)

        memo[pos] = res
        return res
    }

    return dfs(0)
}
