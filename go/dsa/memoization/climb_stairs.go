package memoization

/**
 * 70. Climbing Stairs
 */

/*
 * Top-down approach recursive
 * Idea: Draw the recursion tree then solve the problem
 */
func climbStairs1(n int) int {
	var dfs func(steps int) int
	dfs = func(steps int) int {
		if steps == 0 {
			return 1
		}
		if steps == 1 {
			return 1
		}
		return dfs(steps-1) + dfs(steps-2)
	}
	return dfs(n)
}

/*
 * Top-down approach recursive with memoization
 * Idea: Use a map to store results of subproblems
 */
func climbStairs2(n int) int {
	memo := make([]int, n+1)
	for i := range memo {
		memo[i] = -1
	}

	var dfs func(steps int) int
	dfs = func(steps int) int {
		if steps == 0 {
			return 1
		}
		if steps == 1 {
			return 1
		}
		if memo[steps] != -1 {
			return memo[steps]
		}
		memo[steps] = dfs(steps-1) + dfs(steps-2)
		return memo[steps]
	}
	return dfs(n)
}

/*
 * Bottom-up approach iterative with array
 * Idea: Use an array to mock fuction calls. Start from the base cases till you reach the top initial problem.
 */
func climbStairs3(n int) int {
	memo := make([]int, n+1)
	memo[0], memo[1] = 1, 1
	for i := 2; i <= n; i++ {
		memo[i] = memo[i-2] + memo[i-1]
	}
	return memo[n]
}

/*
 * Bottom-up approach iterative with variables
 * Idea: From the previous solution, we can see that we only need the last two values to calculate the next value.
 * So we can use two variables to store the last two values instead of using an array.
 */
func climbStairs4(n int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return 1
	}
	iMinus1, res := 1, 2
	for i := 2; i < n; i++ {
		res, iMinus1 = iMinus1+res, res
	}
	return res
}
