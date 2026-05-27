import os
import sys
from typing import List

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))


class Solution:
    """
    238. Product of Array Except Self

    # [1, 2, 3, 4]
    # [1, 1, 2, 6]
    # [24,12,4, 1]
    """

    def productExceptSelf(self, nums: List[int]) -> List[int]:
        n = len(nums)
        prefix_product_from_left = [1] * n
        prefix_product_from_right = [1] * n

        for i in range(1, n):
            prefix_product_from_left[i] = prefix_product_from_left[i - 1] * nums[i - 1]

        for i in range(n - 2, -1, -1):
            prefix_product_from_right[i] = (
                prefix_product_from_right[i + 1] * nums[i + 1]
            )

        result = []
        for i in range(0, n):
            result.append(prefix_product_from_left[i] * prefix_product_from_right[i])
        return result

    def product_except_self(self, nums: List[int]) -> List[int]:
        n = len(nums)
        prefix_product = [1] * n

        for i in range(1, n):
            prefix_product[i] = prefix_product[i - 1] * nums[i - 1]

        right = 1
        for i in range(n - 2, -1, -1):
            right *= nums[i + 1]
            prefix_product[i] *= right

        return prefix_product
