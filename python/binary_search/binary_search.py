from typing import List

"""
Am I looking for an exact value?
    YES → Template 1 (left <= right), return mid or -1
    NO  → Template 2 (left < right), finding a boundary
            │
            ├─ First position where condition holds?
            │       right = mid when true → return left
            │
            └─ Last position where condition holds?
                    left = mid+1 when true → return left-1
"""

# The isBadVersion API is already defined for you.
def isBadVersion(version: int) -> bool:
    pass

class Solution:
    """
    33. Search in Rotated Sorted Array

    search for exact match
    """
    def search(self, nums: List[int], target: int) -> int:
        if not nums:
            return -1

        left, right = 0, len(nums)-1
        while left <= right:
            mid = left + (right - left)//2
            if nums[mid] == target:
                return mid

            if nums[mid] < nums[right]: # right half sorted
                if nums[mid] < target <= nums[right]:
                    left = mid + 1
                else:
                    right = mid - 1
            else: # left half sorted
                if nums[left] <= target < nums[mid]:
                    right = mid - 1
                else:
                    left = mid + 1
        return -1


    """
    278. First Bad Version

    search for a boundary
    """
    def firstBadVersion(self, n: int) -> int:
        left, right = 1, n
        while left < right:
            mid = left + (right - left)//2
            if isBadVersion(mid):
                right = mid
            else:
                left = mid + 1
        return left