from typing import List

"""
subsequence => non-contiguous
subarray => contiguous
"""


class Solution:
    """
    1143. Longest Common Subsequence
    """

    """
    300. Longest Increasing Subsequence / LIS

    dp[i] = length of the longest increasing subsequence ending at index i 
    
    Transition:
        dp[i] = max(dp[j] + 1)  for all j < i where nums[j] < nums[i]
        dp[i] = 1               if no such j exists
    
    Base case: dp[i] = 1 for all i (every element is an LIS of length 1 by itself)
    Answer: max(dp)
    """

    def lengthOfLIS(self, nums: List[int]) -> int:
        n = len(nums)
        dp = [1] * n

        res = 1
        for i in range(0, n):
            for j in range(0, i):
                if nums[j] < nums[i]:
                    dp[i] = max(dp[i], dp[j] + 1)
                    res = max(res, dp[i])

        return res

    """
    compute all LIS
    input: [1, 3, 5, 4, 7]
    output: [[1, 3, 5, 7], [1, 3, 4, 7]]

    O(n² × k) time and O(n × k) space where k is the number of LIS (can be exponential)
    """

    def computeLIS(self, nums: List[int]) -> List[List[int]]:
        n = len(nums)
        dp: List[List[int]] = [[[nums[i]]] for i in range(n)]

        for i in range(n):
            max_len = 1
            for j in range(i):
                if nums[j] < nums[i]:
                    for seq in dp[j]:
                        max_len = max(max_len, len(seq) + 1)

            if max_len > 1:
                dp[i] = []
                for j in range(i):
                    if nums[j] < nums[i]:
                        for seq in dp[j]:
                            if len(seq) + 1 == max_len:
                                dp[i].append(seq + [nums[i]])

        max_global = max(len(dp[i][0]) for i in range(n))
        return [seq for i in range(n) for seq in dp[i] if len(seq) == max_global]

    """
    Count the number of LCS
    Complexity: O(n²) time, O(n) space
    """

    def countLIS(self, nums: List[int]) -> int:
        n = len(nums)
        length = [1] * n  # length of LIS ending at i
        count = [1] * n  # number of such LIS

        for i in range(n):
            for j in range(i):
                if nums[j] < nums[i]:
                    if length[j] + 1 > length[i]:
                        length[i] = length[j] + 1
                        count[i] = count[j]
                    elif length[j] + 1 == length[i]:
                        count[i] += count[j]

        max_len = max(length)
        return sum(c for l, c in zip(length, count) if l == max_len)
