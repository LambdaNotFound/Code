from collections import deque


def bfs(start, graph):
    """Generic BFS over an adjacency structure."""
    queue = deque([start])
    visited = {start}  # mark on ENQUEUE, not dequeue, to avoid dupes

    while queue:
        node = queue.popleft()
        for neighbor in graph[node]:
            if neighbor not in visited:
                visited.add(neighbor)
                queue.append(neighbor)
    return visited


def bfs_levels(start, target, graph):
    queue = deque([start])
    visited = {start}
    steps = 0

    while queue:
        for _ in range(len(queue)):  # drain exactly one level
            node = queue.popleft()
            if node == target:
                return steps
            for neighbor in graph[node]:
                if neighbor not in visited:
                    visited.add(neighbor)
                    queue.append(neighbor)
        steps += 1
    return -1


def bfs_grid(grid, start):
    rows, cols = len(grid), len(grid[0])
    queue = deque([start])
    visited = {start}
    DIRS = [(-1, 0), (1, 0), (0, -1), (0, 1)]  # 4-directional

    while queue:
        r, c = queue.popleft()
        for dr, dc in DIRS:
            nr, nc = r + dr, c + dc
            if (
                0 <= nr < rows
                and 0 <= nc < cols
                and (nr, nc) not in visited
                and grid[nr][nc] != obstacle
            ):  # adapt condition
                visited.add((nr, nc))
                queue.append((nr, nc))
