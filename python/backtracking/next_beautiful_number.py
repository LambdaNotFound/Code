from bisect import bisect_right
from collections import Counter


class Solution:
    """
    2048. Next Greater Numerically Balanced Number

    Input: n = 1000
    Output: 1333
    """

    def nextBeautifulNumber(self, n: int) -> int:
        for i in range(n + 1, 1224445):
            count = Counter(str(i))
            if all(count[d] == int(d) for d in count):
                return i

    def next_beautiful_number(self, n: int) -> int:
        def generate(num: int, count: list[int], nums: list[int]) -> None:
            if num > 0 and is_beautiful(count):
                nums.append(num)
            if num > 1224444:
                return

            for d in range(1, 8):
                if count[d] < d:
                    count[d] += 1
                    generate(num * 10 + d, count, nums)
                    count[d] -= 1

        def is_beautiful(count: list[int]) -> bool:
            for d in range(1, 8):
                if count[d] != 0 and count[d] != d:
                    return False
            return True

        nums = []
        generate(0, [0] * 10, nums)
        nums.sort()
        res = bisect_right(nums, n)
        return nums[res]
