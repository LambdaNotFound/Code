
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
    278. First Bad Version
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