
class Solution:
    """
    279. Perfect Squares

    dp[i] = min over k where k*k <= i  of  ( dp[i - k*k] + 1 )
    """
    def numSquares(self, n: int) -> int:
        # Base case: dp[0] = 0 (zero squares needed to make zero)
        # multiplying a list by an integer repeats the list that many times
        dp = [0] + [float('inf')] * n

        for i in range(1, n + 1):
            # Try every perfect square k*k that fits within i.
            # k ranges from 1 to floor(sqrt(i)), so the inner loop is O(sqrt(i)).
            k = 1
            while k * k <= i:
                # Subtract the square k*k and reuse the already-computed subproblem.
                # +1 accounts for using k*k itself as one of the squares.
                dp[i] = min(dp[i], dp[i - k * k] + 1)
                k += 1

        return dp[n]
    
    def numSquaresWithList(n: int) -> list[int]:
        # dp[i] = the actual list of perfect squares summing to i (optimally)
        # dp[0] = [] (empty list — zero squares sum to zero)
        dp = [[] for _ in range(n + 1)]

        for i in range(1, n + 1):
            best = None
            k = 1
            while k * k <= i:
                # Candidate decomposition: take the optimal list for (i - k*k),
                # then append k*k as one more square.
                candidate = dp[i - k * k] + [k * k]
                if best is None or len(candidate) < len(best):
                    best = candidate
                k += 1
            dp[i] = best

        return dp[n]