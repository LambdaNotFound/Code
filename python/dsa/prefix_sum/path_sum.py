from collections import defaultdict
import sys, os
from typing import Optional
sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from common import TreeNode

class Solution:
    """
    437. Path Sum III

    prefixSum[j] - prefixSum[i] = k => prefixSum[i] = prefixSum[j] - k
    """
    def pathSum(self, root: Optional[TreeNode], targetSum: int) -> int:
        prefix_counts = defaultdict(int)
        prefix_counts[0] = 1   # empty path base case

        def dfs(node, curr_sum) -> int:
            if not node:
                return 0

            curr_sum += node.val

            # how many paths ending here sum to target
            count = prefix_counts[curr_sum - targetSum]

            # explore children with updated prefix
            prefix_counts[curr_sum] += 1
            count += dfs(node.left, curr_sum)
            count += dfs(node.right, curr_sum)
            prefix_counts[curr_sum] -= 1   # backtrack

            return count

        return dfs(root, 0)