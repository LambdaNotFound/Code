import bisect


class Solution(object):
    """
    475. Heaters

    """

    def findRadius(houses, heaters):
        heaters.sort()

        max_min_dist = 0
        for house in houses:
            # sort.Search with house < heaters[i] = first index > house = bisect_right
            right_idx = bisect.bisect_right(heaters, house)
            left_idx = right_idx - 1

            min_dist = float("inf")
            if right_idx < len(heaters):
                min_dist = min(min_dist, heaters[right_idx] - house)
            if left_idx >= 0:
                min_dist = min(min_dist, house - heaters[left_idx])

            max_min_dist = max(max_min_dist, min_dist)

        return max_min_dist
