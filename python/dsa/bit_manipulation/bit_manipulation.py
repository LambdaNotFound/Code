class Solution:
    """
    268. Missing Number
    """

    def missingNumber(self, nums: list[int]) -> int:
        nums.sort()
        for i, num in enumerate(nums):
            if num != i:
                return i
        return len(nums)
