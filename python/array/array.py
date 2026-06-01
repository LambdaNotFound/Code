from typing import List


class Solution:
    """
    You are given two arrays of integers, heights1 and heights2, representing the heights of players from two different teams. Your task is to determine if it's possible to arrange the teams for a photograph according to the following rules:

    One team must stand entirely in front of the other.
    Every player in the front row must be strictly shorter than the player standing directly behind them.
    Players can be rearranged within their own team's row.
    The team with fewer players must stand in the front row. If the teams are of equal size, either team can be the front row.
    Your function should return True if such an arrangement is possible, and False otherwise.
    """

    def can_arrange(self, heights1, heights2):
        def can_match(front, back):
            """Can every player in `front` stand before a distinct, strictly taller `back` player?"""
            front = sorted(front)
            back = sorted(back)
            i = 0  # index of the shortest still-unmatched front player
            for h in back:
                if i < len(front) and h > front[i]:
                    i += 1  # this back player covers front[i]
            return i == len(front)  # did everyone in front get matched?

        if len(heights1) == len(heights2):
            # equal size — either team may lead, so try both
            return can_match(heights1, heights2) or can_match(heights2, heights1)
        if len(heights1) < len(heights2):
            return can_match(heights1, heights2)
        return can_match(heights2, heights1)

    """
    41. First Missing Positive

    goal: [1, 2, 3, 4, 5]
        0  1  2  3  4
    """

    def firstMissingPositive(self, nums: List[int]) -> int:
        n = len(nums)

        # Phase 1: cyclic sort — place value v at index v-1 for v in [1, n]
        for i, _ in enumerate(nums):
            while 0 <= nums[i] - 1 <= n - 1 and nums[i] != nums[nums[i] - 1]:
                target = nums[i] - 1  # cache before swap — see notes below
                nums[i], nums[target] = nums[target], nums[i]

        # Phase 2: find first mismatch
        for i, _ in enumerate(nums):
            if nums[i] != i + 1:
                return i + 1

        return n + 1  # all [1..n] present → answer is n+1
