from collections import deque
from typing import Optional, List
from common import TreeNode

class Solution:
    """
    199. Binary Tree Right Side View
    """
    def rightSideView(self, root: Optional[TreeNode]) -> List[int]:
        if root is None:
            return []
        
        queue: deque[TreeNode] = deque()
        queue.append(root)

        result: List[int] = []
        while queue:
            result.append(queue[-1].val)

            count = len(queue)
            for _ in range(count):
                node = queue.popleft()
                if node.left:
                    queue.append(node.left)
                if node.right:
                    queue.append(node.right)
        return result
